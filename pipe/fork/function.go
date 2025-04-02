//
// Copyright (C) 2022 - 2025 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package fork

import (
	"context"

	"github.com/fogfish/golem/pipe/v2"
)

// Pure effect over category A ⟼ B
type E[A, B any] = func(A) B

// Either effect over category A ⟼ (B, error)
type EitherE[A, B any] = func(A) (B, error)

// Arrow over functor 𝓕: A ⟼ B
type Arrow[A, B any] = func(context.Context, A, chan<- B) error

// Go channel morphism 𝑓: A ⟼ B
type F[A, B any] interface {
	Apply(A) (B, error)
	errch(cap int) chan error
	catch(context.Context, error, chan<- error) bool
	pipef() pipe.F[A, B]
}

// Lift pure effect into morphism 𝑓: A ⟼ B
func Pure[A, B any](f E[A, B]) F[A, B] {
	ff := func(a A) (B, error) { return f(a), nil }
	return pure[A, B](ff)
}

// Lift either effect into morphism 𝑓: A ⟼ B
// The failure of morphism causes the failure of channel, aborts the computation.
func Lift[A, B any](f EitherE[A, B]) F[A, B] {
	return pure[A, B](f)
}

type pure[A, B any] EitherE[A, B]

func (f pure[A, B]) Apply(a A) (B, error) {
	return EitherE[A, B](f)(a)
}

//lint:ignore U1000 false positive
func (f pure[A, B]) errch(_ int) chan error {
	return make(chan error, 1)
}

//lint:ignore U1000 false positive
func (f pure[A, B]) catch(ctx context.Context, err error, exx chan<- error) bool {
	exx <- err
	return false
}

//lint:ignore U1000 false positive
func (f pure[A, B]) pipef() pipe.F[A, B] {
	return pipe.Lift(f)
}

// Lift either effect into morphism 𝑓: A ⟼ B
// The failure of morphism causes the failure of step, continues the computation.
func Try[A, B any](f EitherE[A, B]) F[A, B] {
	return try[A, B](f)
}

type try[A, B any] EitherE[A, B]

func (f try[A, B]) Apply(a A) (B, error) {
	return EitherE[A, B](f)(a)
}

//lint:ignore U1000 false positive
func (f try[A, B]) errch(cap int) chan error {
	return make(chan error, cap)
}

//lint:ignore U1000 false positive
func (f try[A, B]) catch(ctx context.Context, err error, exx chan<- error) bool {
	select {
	case exx <- err:
	case <-ctx.Done():
		return false
	}
	return true
}

//lint:ignore U1000 false positive
func (f try[A, B]) pipef() pipe.F[A, B] {
	return pipe.Try(f)
}

//------------------------------------------------------------------------------

// Go channel functor 𝓕: A ⟼ B
type FF[A, B any] interface {
	Apply(context.Context, A, chan<- B) error
	errch(cap int) chan error
	catch(context.Context, error, chan<- error) bool
}

// Lift arrow into functor morphism 𝓕: A ⟼ B
// The failure of morphism causes the failure of channel, aborts the computation.
func LiftF[A, B any](f Arrow[A, B]) FF[A, B] {
	return puref[A, B](f)
}

type puref[A, B any] Arrow[A, B]

func (f puref[A, B]) Apply(ctx context.Context, a A, b chan<- B) error {
	return Arrow[A, B](f)(ctx, a, b)
}

//lint:ignore U1000 false positive
func (f puref[A, B]) errch(_ int) chan error {
	return make(chan error, 1)
}

//lint:ignore U1000 false positive
func (f puref[A, B]) catch(ctx context.Context, err error, exx chan<- error) bool {
	exx <- err
	return false
}

// Lift arrow into functor morphism 𝓕: A ⟼ B
// The failure of morphism causes the failure of step, continues the computation.
func TryF[A, B any](f Arrow[A, B]) FF[A, B] {
	return tryf[A, B](f)
}

type tryf[A, B any] Arrow[A, B]

func (f tryf[A, B]) Apply(ctx context.Context, a A, b chan<- B) error {
	return Arrow[A, B](f)(ctx, a, b)
}

//lint:ignore U1000 false positive
func (f tryf[A, B]) errch(cap int) chan error {
	return make(chan error, cap)
}

//lint:ignore U1000 false positive
func (f tryf[A, B]) catch(ctx context.Context, err error, exx chan<- error) bool {
	select {
	case exx <- err:
	case <-ctx.Done():
		return false
	}
	return true
}
