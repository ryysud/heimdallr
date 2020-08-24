// Generated by protoc-ddl.
// protoc-gen-entity: v0.1
package entity

import (
	"bytes"
	"sync"
	"time"

	"go.f110.dev/protoc-ddl"
)

var _ = time.Time{}
var _ = bytes.Buffer{}

type Column struct {
	Name  string
	Value interface{}
}

type User struct {
	Id        int32
	Identity  string
	Admin     bool
	Type      string
	Comment   string
	CreatedAt time.Time
	UpdatedAt *time.Time

	mu   sync.Mutex
	mark *User
}

func (e *User) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *User) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Identity != e.mark.Identity ||
		e.Admin != e.mark.Admin ||
		e.Type != e.mark.Type ||
		e.Comment != e.mark.Comment ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *User) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Identity != e.mark.Identity {
		res = append(res, ddl.Column{Name: "identity", Value: e.Identity})
	}
	if e.Admin != e.mark.Admin {
		res = append(res, ddl.Column{Name: "admin", Value: e.Admin})
	}
	if e.Type != e.mark.Type {
		res = append(res, ddl.Column{Name: "type", Value: e.Type})
	}
	if e.Comment != e.mark.Comment {
		res = append(res, ddl.Column{Name: "comment", Value: e.Comment})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *User) Copy() *User {
	n := &User{
		Id:        e.Id,
		Identity:  e.Identity,
		Admin:     e.Admin,
		Type:      e.Type,
		Comment:   e.Comment,
		CreatedAt: e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	return n
}

type UserState struct {
	Id        int32
	State     string
	Unique    string
	CreatedAt time.Time
	UpdatedAt *time.Time

	mu   sync.Mutex
	mark *UserState
}

func (e *UserState) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *UserState) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.State != e.mark.State ||
		e.Unique != e.mark.Unique ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *UserState) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.State != e.mark.State {
		res = append(res, ddl.Column{Name: "state", Value: e.State})
	}
	if e.Unique != e.mark.Unique {
		res = append(res, ddl.Column{Name: "unique", Value: e.Unique})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *UserState) Copy() *UserState {
	n := &UserState{
		Id:        e.Id,
		State:     e.State,
		Unique:    e.Unique,
		CreatedAt: e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	return n
}

type RoleBinding struct {
	Id         int32
	UserId     int32
	Role       string
	Maintainer bool
	CreatedAt  time.Time
	UpdatedAt  *time.Time

	User *User

	mu   sync.Mutex
	mark *RoleBinding
}

func (e *RoleBinding) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *RoleBinding) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.UserId != e.mark.UserId ||
		e.Role != e.mark.Role ||
		e.Maintainer != e.mark.Maintainer ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *RoleBinding) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if e.Role != e.mark.Role {
		res = append(res, ddl.Column{Name: "role", Value: e.Role})
	}
	if e.Maintainer != e.mark.Maintainer {
		res = append(res, ddl.Column{Name: "maintainer", Value: e.Maintainer})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *RoleBinding) Copy() *RoleBinding {
	n := &RoleBinding{
		Id:         e.Id,
		UserId:     e.UserId,
		Role:       e.Role,
		Maintainer: e.Maintainer,
		CreatedAt:  e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type AccessToken struct {
	Id        int32
	Name      string
	Value     string
	UserId    int32
	IssuerId  int32
	CreatedAt time.Time
	UpdatedAt *time.Time

	User   *User
	Issuer *User

	mu   sync.Mutex
	mark *AccessToken
}

func (e *AccessToken) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *AccessToken) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Name != e.mark.Name ||
		e.Value != e.mark.Value ||
		e.UserId != e.mark.UserId ||
		e.IssuerId != e.mark.IssuerId ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *AccessToken) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Name != e.mark.Name {
		res = append(res, ddl.Column{Name: "name", Value: e.Name})
	}
	if e.Value != e.mark.Value {
		res = append(res, ddl.Column{Name: "value", Value: e.Value})
	}
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if e.IssuerId != e.mark.IssuerId {
		res = append(res, ddl.Column{Name: "issuer_id", Value: e.IssuerId})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *AccessToken) Copy() *AccessToken {
	n := &AccessToken{
		Id:        e.Id,
		Name:      e.Name,
		Value:     e.Value,
		UserId:    e.UserId,
		IssuerId:  e.IssuerId,
		CreatedAt: e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	if e.Issuer != nil {
		n.Issuer = e.Issuer.Copy()
	}
	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type Token struct {
	Id       int32
	Token    string
	UserId   int32
	IssuedAt time.Time

	User *User

	mu   sync.Mutex
	mark *Token
}

func (e *Token) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Token) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Token != e.mark.Token ||
		e.UserId != e.mark.UserId ||
		!e.IssuedAt.Equal(e.mark.IssuedAt)
}

func (e *Token) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Token != e.mark.Token {
		res = append(res, ddl.Column{Name: "token", Value: e.Token})
	}
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if !e.IssuedAt.Equal(e.mark.IssuedAt) {
		res = append(res, ddl.Column{Name: "issued_at", Value: e.IssuedAt})
	}

	return res
}

func (e *Token) Copy() *Token {
	n := &Token{
		Id:       e.Id,
		Token:    e.Token,
		UserId:   e.UserId,
		IssuedAt: e.IssuedAt,
	}

	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type Code struct {
	Id              int32
	Code            string
	Challenge       string
	ChallengeMethod string
	UserId          int32
	IssuedAt        time.Time

	User *User

	mu   sync.Mutex
	mark *Code
}

func (e *Code) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Code) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Code != e.mark.Code ||
		e.Challenge != e.mark.Challenge ||
		e.ChallengeMethod != e.mark.ChallengeMethod ||
		e.UserId != e.mark.UserId ||
		!e.IssuedAt.Equal(e.mark.IssuedAt)
}

func (e *Code) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Code != e.mark.Code {
		res = append(res, ddl.Column{Name: "code", Value: e.Code})
	}
	if e.Challenge != e.mark.Challenge {
		res = append(res, ddl.Column{Name: "challenge", Value: e.Challenge})
	}
	if e.ChallengeMethod != e.mark.ChallengeMethod {
		res = append(res, ddl.Column{Name: "challenge_method", Value: e.ChallengeMethod})
	}
	if e.UserId != e.mark.UserId {
		res = append(res, ddl.Column{Name: "user_id", Value: e.UserId})
	}
	if !e.IssuedAt.Equal(e.mark.IssuedAt) {
		res = append(res, ddl.Column{Name: "issued_at", Value: e.IssuedAt})
	}

	return res
}

func (e *Code) Copy() *Code {
	n := &Code{
		Id:              e.Id,
		Code:            e.Code,
		Challenge:       e.Challenge,
		ChallengeMethod: e.ChallengeMethod,
		UserId:          e.UserId,
		IssuedAt:        e.IssuedAt,
	}

	if e.User != nil {
		n.User = e.User.Copy()
	}

	return n
}

type Relay struct {
	Id          int32
	Name        string
	Addr        string
	FromAddr    string
	ConnectedAt time.Time
	CreatedAt   time.Time
	UpdatedAt   *time.Time

	mu   sync.Mutex
	mark *Relay
}

func (e *Relay) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Relay) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Name != e.mark.Name ||
		e.Addr != e.mark.Addr ||
		e.FromAddr != e.mark.FromAddr ||
		!e.ConnectedAt.Equal(e.mark.ConnectedAt) ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *Relay) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Name != e.mark.Name {
		res = append(res, ddl.Column{Name: "name", Value: e.Name})
	}
	if e.Addr != e.mark.Addr {
		res = append(res, ddl.Column{Name: "addr", Value: e.Addr})
	}
	if e.FromAddr != e.mark.FromAddr {
		res = append(res, ddl.Column{Name: "from_addr", Value: e.FromAddr})
	}
	if !e.ConnectedAt.Equal(e.mark.ConnectedAt) {
		res = append(res, ddl.Column{Name: "connected_at", Value: e.ConnectedAt})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *Relay) Copy() *Relay {
	n := &Relay{
		Id:          e.Id,
		Name:        e.Name,
		Addr:        e.Addr,
		FromAddr:    e.FromAddr,
		ConnectedAt: e.ConnectedAt,
		CreatedAt:   e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	return n
}

type SerialNumber struct {
	Id           int64
	SerialNumber []byte

	mu   sync.Mutex
	mark *SerialNumber
}

func (e *SerialNumber) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *SerialNumber) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return !bytes.Equal(e.SerialNumber, e.mark.SerialNumber)
}

func (e *SerialNumber) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if !bytes.Equal(e.SerialNumber, e.mark.SerialNumber) {
		res = append(res, ddl.Column{Name: "serial_number", Value: e.SerialNumber})
	}

	return res
}

func (e *SerialNumber) Copy() *SerialNumber {
	n := &SerialNumber{
		Id:           e.Id,
		SerialNumber: e.SerialNumber,
	}

	return n
}

type SignedCertificate struct {
	Id             int32
	Certificate    []byte
	SerialNumberId int64
	P12            []byte
	Agent          bool
	Comment        string
	IssuedAt       time.Time

	SerialNumber *SerialNumber

	mu   sync.Mutex
	mark *SignedCertificate
}

func (e *SignedCertificate) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *SignedCertificate) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return !bytes.Equal(e.Certificate, e.mark.Certificate) ||
		e.SerialNumberId != e.mark.SerialNumberId ||
		!bytes.Equal(e.P12, e.mark.P12) ||
		e.Agent != e.mark.Agent ||
		e.Comment != e.mark.Comment ||
		!e.IssuedAt.Equal(e.mark.IssuedAt)
}

func (e *SignedCertificate) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if !bytes.Equal(e.Certificate, e.mark.Certificate) {
		res = append(res, ddl.Column{Name: "certificate", Value: e.Certificate})
	}
	if e.SerialNumberId != e.mark.SerialNumberId {
		res = append(res, ddl.Column{Name: "serial_number_id", Value: e.SerialNumberId})
	}
	if !bytes.Equal(e.P12, e.mark.P12) {
		res = append(res, ddl.Column{Name: "p12", Value: e.P12})
	}
	if e.Agent != e.mark.Agent {
		res = append(res, ddl.Column{Name: "agent", Value: e.Agent})
	}
	if e.Comment != e.mark.Comment {
		res = append(res, ddl.Column{Name: "comment", Value: e.Comment})
	}
	if !e.IssuedAt.Equal(e.mark.IssuedAt) {
		res = append(res, ddl.Column{Name: "issued_at", Value: e.IssuedAt})
	}

	return res
}

func (e *SignedCertificate) Copy() *SignedCertificate {
	n := &SignedCertificate{
		Id:             e.Id,
		Certificate:    e.Certificate,
		SerialNumberId: e.SerialNumberId,
		P12:            e.P12,
		Agent:          e.Agent,
		Comment:        e.Comment,
		IssuedAt:       e.IssuedAt,
	}

	if e.SerialNumber != nil {
		n.SerialNumber = e.SerialNumber.Copy()
	}

	return n
}

type RevokedCertificate struct {
	Id           int32
	CommonName   string
	SerialNumber []byte
	Agent        bool
	Comment      string
	RevokedAt    time.Time
	IssuedAt     time.Time
	CreatedAt    time.Time
	UpdatedAt    *time.Time

	mu   sync.Mutex
	mark *RevokedCertificate
}

func (e *RevokedCertificate) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *RevokedCertificate) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.CommonName != e.mark.CommonName ||
		!bytes.Equal(e.SerialNumber, e.mark.SerialNumber) ||
		e.Agent != e.mark.Agent ||
		e.Comment != e.mark.Comment ||
		!e.RevokedAt.Equal(e.mark.RevokedAt) ||
		!e.IssuedAt.Equal(e.mark.IssuedAt) ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *RevokedCertificate) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.CommonName != e.mark.CommonName {
		res = append(res, ddl.Column{Name: "common_name", Value: e.CommonName})
	}
	if !bytes.Equal(e.SerialNumber, e.mark.SerialNumber) {
		res = append(res, ddl.Column{Name: "serial_number", Value: e.SerialNumber})
	}
	if e.Agent != e.mark.Agent {
		res = append(res, ddl.Column{Name: "agent", Value: e.Agent})
	}
	if e.Comment != e.mark.Comment {
		res = append(res, ddl.Column{Name: "comment", Value: e.Comment})
	}
	if !e.RevokedAt.Equal(e.mark.RevokedAt) {
		res = append(res, ddl.Column{Name: "revoked_at", Value: e.RevokedAt})
	}
	if !e.IssuedAt.Equal(e.mark.IssuedAt) {
		res = append(res, ddl.Column{Name: "issued_at", Value: e.IssuedAt})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *RevokedCertificate) Copy() *RevokedCertificate {
	n := &RevokedCertificate{
		Id:           e.Id,
		CommonName:   e.CommonName,
		SerialNumber: e.SerialNumber,
		Agent:        e.Agent,
		Comment:      e.Comment,
		RevokedAt:    e.RevokedAt,
		IssuedAt:     e.IssuedAt,
		CreatedAt:    e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	return n
}

type Node struct {
	Id        int32
	Hostname  string
	CreatedAt time.Time
	UpdatedAt *time.Time

	mu   sync.Mutex
	mark *Node
}

func (e *Node) ResetMark() {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.mark = e.Copy()
}

func (e *Node) IsChanged() bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	return e.Hostname != e.mark.Hostname ||
		!e.CreatedAt.Equal(e.mark.CreatedAt) ||
		((e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil))
}

func (e *Node) ChangedColumn() []ddl.Column {
	e.mu.Lock()
	defer e.mu.Unlock()

	res := make([]ddl.Column, 0)
	if e.Hostname != e.mark.Hostname {
		res = append(res, ddl.Column{Name: "hostname", Value: e.Hostname})
	}
	if !e.CreatedAt.Equal(e.mark.CreatedAt) {
		res = append(res, ddl.Column{Name: "created_at", Value: e.CreatedAt})
	}
	if (e.UpdatedAt != nil && (e.mark.UpdatedAt == nil || !e.UpdatedAt.Equal(*e.mark.UpdatedAt))) || (e.UpdatedAt == nil && e.mark.UpdatedAt != nil) {
		if e.UpdatedAt != nil {
			res = append(res, ddl.Column{Name: "updated_at", Value: *e.UpdatedAt})
		} else {
			res = append(res, ddl.Column{Name: "updated_at", Value: nil})
		}
	}

	return res
}

func (e *Node) Copy() *Node {
	n := &Node{
		Id:        e.Id,
		Hostname:  e.Hostname,
		CreatedAt: e.CreatedAt,
	}
	if e.UpdatedAt != nil {
		v := *e.UpdatedAt
		n.UpdatedAt = &v
	}

	return n
}
