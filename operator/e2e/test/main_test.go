package test

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"testing"
	"time"

	_ "github.com/smartystreets/goconvey/convey"
	"golang.org/x/xerrors"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"

	"go.f110.dev/heimdallr/operator/e2e/e2eutil"
	"go.f110.dev/heimdallr/operator/e2e/framework"
	clientset "go.f110.dev/heimdallr/operator/pkg/client/versioned"
	"go.f110.dev/heimdallr/operator/pkg/controllers"
	informers "go.f110.dev/heimdallr/operator/pkg/informers/externalversions"
	"go.f110.dev/heimdallr/pkg/config/configv2"
	"go.f110.dev/heimdallr/pkg/k8s"
	"go.f110.dev/heimdallr/pkg/k8s/kind"
	"go.f110.dev/heimdallr/pkg/logger"
)

func init() {
	framework.Flags(flag.CommandLine)
}

func setupSuite(id string) (*kind.Cluster, error) {
	k8sCluster, err := kind.NewCluster(framework.Config.KindFile, "e2e-"+id, "")
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if err := k8sCluster.Create(framework.Config.ClusterVersion, 3); err != nil {
		log.Fatalf("Could not create a cluster: %v", err)
	}
	kubeConfig := k8sCluster.KubeConfig()
	log.Printf("KubeConfig: %s", kubeConfig)

	cfg, err := k8sCluster.RESTConfig()
	if err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	RESTConfig = cfg

	if err := k8sCluster.WaitReady(context.TODO()); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if framework.Config.ProxyImageFile != "" || framework.Config.RPCImageFile != "" || framework.Config.DashboardImageFile != "" {
		images := []*kind.ContainerImageFile{
			{
				File:       framework.Config.ProxyImageFile,
				Repository: controllers.ProxyImageRepository,
				Tag:        "e2e",
			},
			{
				File:       framework.Config.RPCImageFile,
				Repository: controllers.RPCServerImageRepository,
				Tag:        "e2e",
			},
			{
				File:       framework.Config.DashboardImageFile,
				Repository: controllers.DashboardImageRepository,
				Tag:        "e2e",
			},
		}
		if err := k8sCluster.LoadImageFiles(images...); err != nil {
			log.Fatal(err)
		}

		framework.Config.ProxyVersion = "e2e"
	}

	if err := kind.InstallCertManager(cfg, "operator-e2e"); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}
	if err := kind.InstallMinIO(cfg, "operator-e2e"); err != nil {
		return nil, xerrors.Errorf(": %w", err)
	}

	return k8sCluster, nil
}

func TestMain(m *testing.M) {
	flag.Parse()

	rand.Seed(framework.Config.RandomSeed)
	id := e2eutil.MakeId()

	logLevel := "fatal"
	if framework.Config.Verbose {
		logLevel = "debug"
	}
	if err := logger.OverrideKlog(&configv2.Logger{Level: logLevel}); err != nil {
		panic(err)
		return
	}
	fs := flag.NewFlagSet("e2e", flag.ContinueOnError)
	klog.InitFlags(fs)
	if err := fs.Parse([]string{"-stderrthreshold=FATAL", fmt.Sprintf("-logtostderr=%v", framework.Config.Verbose)}); err != nil {
		log.Fatal(err)
	}
	klog.SetOutput(ioutil.Discard)

	log.Printf("%+v", framework.Config)

	var k8sCluster *kind.Cluster
	framework.BeforeSuite(func() {
		crd, err := k8s.ReadCRDFile(framework.Config.CRDDir)
		if err != nil {
			log.Fatal(err)
		}

		if v, err := setupSuite(id); err != nil {
			log.Fatal(err)
		} else {
			k8sCluster = v
		}

		// Create CustomResourceDefinition
		cfg, err := k8sCluster.RESTConfig()
		if err != nil {
			log.Fatal(err)
		}
		if err := k8s.EnsureCRD(cfg, crd, 3*time.Minute); err != nil {
			log.Fatal(err)
		}
		kubeClient, err := k8sCluster.Clientset()
		if err != nil {
			log.Fatal(err)
		}
		proxyClient, err := clientset.NewForConfig(cfg)
		if err != nil {
			log.Fatal(err)
		}

		lock, err := resourcelock.New(
			resourcelock.LeasesResourceLock,
			"default",
			"e2e",
			kubeClient.CoreV1(),
			kubeClient.CoordinationV1(),
			resourcelock.ResourceLockConfig{Identity: id},
		)
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			leaderelection.RunOrDie(context.TODO(), leaderelection.LeaderElectionConfig{
				Lock:            lock,
				ReleaseOnCancel: true,
				LeaseDuration:   30 * time.Second,
				RenewDeadline:   15 * time.Second,
				RetryPeriod:     5 * time.Second,
				Callbacks: leaderelection.LeaderCallbacks{
					OnStartedLeading: func(ctx context.Context) {
						coreSharedInformerFactory := kubeinformers.NewSharedInformerFactory(kubeClient, 30*time.Second)
						sharedInformerFactory := informers.NewSharedInformerFactory(proxyClient, 30*time.Second)

						e, err := controllers.NewEtcdController(
							sharedInformerFactory,
							coreSharedInformerFactory,
							kubeClient,
							proxyClient,
							cfg,
							"cluster.local",
							true,
							http.DefaultTransport,
							nil,
						)
						if err != nil {
							log.Fatal(err)
						}

						c, err := controllers.NewProxyController(sharedInformerFactory, coreSharedInformerFactory, kubeClient, proxyClient)
						if err != nil {
							log.Fatal(err)
						}

						sharedInformerFactory.Start(ctx.Done())
						coreSharedInformerFactory.Start(ctx.Done())

						var wg sync.WaitGroup
						wg.Add(1)
						go func() {
							defer wg.Done()

							e.Run(ctx, 1)
						}()

						wg.Add(1)
						go func() {
							defer wg.Done()

							c.Run(ctx, 1)
						}()

						wg.Wait()
					},
					OnStoppedLeading: func() {},
					OnNewLeader: func(identity string) {
						if identity == id {
							return
						}
					},
				},
			})
		}()
	})

	framework.AfterSuite(func() {
		if k8sCluster != nil && !framework.Config.Retain {
			if err := k8sCluster.Delete(); err != nil {
				log.Fatalf("Could not delete a cluster: %v", err)
			}
		}
	})

	framework.RunSpec(m)
}
