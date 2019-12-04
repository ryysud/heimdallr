load("@io_bazel_rules_go//go:def.bzl", "go_context", "go_rule")
load("@bazel_skylib//lib:shell.bzl", "shell")

def _go_vendor_impl(ctx):
    go = go_context(ctx)
    out = ctx.actions.declare_file(ctx.label.name + ".sh")
    dir = ctx.attr.dir
    substitutions = {
        "@@GO@@": shell.quote(go.go.path),
        "@@GAZELLE@@": shell.quote(ctx.executable._gazelle.short_path),
        "@@DIR@@": shell.quote(dir),
        "@@ARGS@@": shell.array_literal(ctx.attr.extra_args),
    }
    ctx.actions.expand_template(
        template = ctx.file._template,
        output = out,
        substitutions = substitutions,
        is_executable = True,
    )
    runfiles = ctx.runfiles(files = [go.go, ctx.executable._gazelle])
    return [
        DefaultInfo(
            runfiles = runfiles,
            executable = out,
        ),
    ]

_go_vendor = go_rule(
    implementation = _go_vendor_impl,
    executable = True,
    attrs = {
        "dir": attr.string(),
        "extra_args": attr.string_list(),
        "_template": attr.label(
            default = "//build/rules/go:vendor.bash",
            allow_single_file = True,
        ),
        "_gazelle": attr.label(
            default = "@bazel_gazelle//cmd/gazelle",
            executable = True,
            cfg = "host",
        ),
    },
)

def go_vendor(name, **kwargs):
    if not "dir" in kwargs:
        dir = native.package_name()
        kwargs["dir"] = dir

    _go_vendor(name = name, **kwargs)