//
// Copyright (C) 2022 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/golem
//

package pure

func Pipe[A, B, C any](
	ab func(A) B,
	bc func(B) C,
) func(A) C {
	return func(a A) C { return bc(ab(a)) }
}

func Pipe3[A, B, C, D any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
) func(A) D {
	return func(a A) D { return cd(bc(ab(a))) }
}

func Pipe4[A, B, C, D, E any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
) func(A) E {
	return func(a A) E { return de(cd(bc(ab(a)))) }
}

func Pipe5[A, B, C, D, E, F any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
) func(A) F {
	return func(a A) F { return ef(de(cd(bc(ab(a))))) }
}

func Pipe6[A, B, C, D, E, F, G any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
) func(A) G {
	return func(a A) G { return fg(ef(de(cd(bc(ab(a)))))) }
}

func Pipe7[A, B, C, D, E, F, G, H any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
) func(A) H {
	return func(a A) H { return gh(fg(ef(de(cd(bc(ab(a))))))) }
}

func Pipe8[A, B, C, D, E, F, G, H, I any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
) func(A) I {
	return func(a A) I { return hi(gh(fg(ef(de(cd(bc(ab(a)))))))) }
}

func Pipe9[A, B, C, D, E, F, G, H, I, J any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
) func(A) J {
	return func(a A) J { return ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))) }
}

func Pipe10[A, B, C, D, E, F, G, H, I, J, K any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
) func(A) K {
	return func(a A) K { return jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))) }
}

func Pipe11[A, B, C, D, E, F, G, H, I, J, K, L any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
) func(A) L {
	return func(a A) L { return kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))))) }
}

func Pipe12[A, B, C, D, E, F, G, H, I, J, K, L, M any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
) func(A) M {
	return func(a A) M { return lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))))) }
}

func Pipe13[A, B, C, D, E, F, G, H, I, J, K, L, M, N any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
) func(A) N {
	return func(a A) N { return mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))))))) }
}

func Pipe14[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
) func(A) O {
	return func(a A) O { return no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))))))) }
}

func Pipe15[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
) func(A) P {
	return func(a A) P { return op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))))))))) }
}

func Pipe16[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
	pq func(P) Q,
) func(A) Q {
	return func(a A) Q { return pq(op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))))))))) }
}

func Pipe17[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
	pq func(P) Q,
	qr func(Q) R,
) func(A) R {
	return func(a A) R { return qr(pq(op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))))))))))) }
}

func Pipe18[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
	pq func(P) Q,
	qr func(Q) R,
	rs func(R) S,
) func(A) S {
	return func(a A) S { return rs(qr(pq(op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))))))))))) }
}

func Pipe19[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
	pq func(P) Q,
	qr func(Q) R,
	rs func(R) S,
	st func(S) T,
) func(A) T {
	return func(a A) T { return st(rs(qr(pq(op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a))))))))))))))))))) }
}

func Pipe20[A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U any](
	ab func(A) B,
	bc func(B) C,
	cd func(C) D,
	de func(D) E,
	ef func(E) F,
	fg func(F) G,
	gh func(G) H,
	hi func(H) I,
	ij func(I) J,
	jk func(J) K,
	kl func(K) L,
	lm func(L) M,
	mn func(M) N,
	no func(N) O,
	op func(O) P,
	pq func(P) Q,
	qr func(Q) R,
	rs func(R) S,
	st func(S) T,
	tu func(T) U,
) func(A) U {
	return func(a A) U { return tu(st(rs(qr(pq(op(no(mn(lm(kl(jk(ij(hi(gh(fg(ef(de(cd(bc(ab(a)))))))))))))))))))) }
}
