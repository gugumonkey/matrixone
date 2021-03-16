// Code generated by command: go run avx2.go -out min/avx2.s -stubs min/avx2_stubs.go. DO NOT EDIT.

#include "textflag.h"

// func int8MinAvx2Asm(x []int8, r []int8)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·int8MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x000000000000007f, BX
	MOVQ         BX, X0
	VPBROADCASTB X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int8MinBlockLoop:
	CMPQ    DX, $0x000000c0
	JL      int8MinTailLoop
	VPMINSB (AX), Y1, Y1
	VPMINSB 32(AX), Y2, Y2
	VPMINSB 64(AX), Y3, Y3
	VPMINSB 96(AX), Y4, Y4
	VPMINSB 128(AX), Y5, Y5
	VPMINSB 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x000000c0, DX
	JMP     int8MinBlockLoop

int8MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      int8MinDone
	VPMINSB (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000020, DX
	JMP     int8MinTailLoop

int8MinDone:
	VPMINSB      Y1, Y2, Y1
	VPMINSB      Y1, Y3, Y1
	VPMINSB      Y1, Y4, Y1
	VPMINSB      Y1, Y5, Y1
	VPMINSB      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINSB       X0, X1
	MOVOU        X1, (CX)
	RET

// func int16MinAvx2Asm(x []int16, r []int16)
// Requires: AVX, AVX2, SSE2
TEXT ·int16MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x0000000000007fff, BX
	MOVQ         BX, X0
	VPBROADCASTW X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int16MinBlockLoop:
	CMPQ    DX, $0x00000060
	JL      int16MinTailLoop
	VPMINSW (AX), Y1, Y1
	VPMINSW 32(AX), Y2, Y2
	VPMINSW 64(AX), Y3, Y3
	VPMINSW 96(AX), Y4, Y4
	VPMINSW 128(AX), Y5, Y5
	VPMINSW 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000060, DX
	JMP     int16MinBlockLoop

int16MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      int16MinDone
	VPMINSW (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000010, DX
	JMP     int16MinTailLoop

int16MinDone:
	VPMINSW      Y1, Y2, Y1
	VPMINSW      Y1, Y3, Y1
	VPMINSW      Y1, Y4, Y1
	VPMINSW      Y1, Y5, Y1
	VPMINSW      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINSW       X0, X1
	MOVOU        X1, (CX)
	RET

// func int32MinAvx2Asm(x []int32, r []int32)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·int32MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x000000007fffffff, BX
	MOVQ         BX, X0
	VPBROADCASTD X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

int32MinBlockLoop:
	CMPQ    DX, $0x00000030
	JL      int32MinTailLoop
	VPMINSD (AX), Y1, Y1
	VPMINSD 32(AX), Y2, Y2
	VPMINSD 64(AX), Y3, Y3
	VPMINSD 96(AX), Y4, Y4
	VPMINSD 128(AX), Y5, Y5
	VPMINSD 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000030, DX
	JMP     int32MinBlockLoop

int32MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      int32MinDone
	VPMINSD (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000008, DX
	JMP     int32MinTailLoop

int32MinDone:
	VPMINSD      Y1, Y2, Y1
	VPMINSD      Y1, Y3, Y1
	VPMINSD      Y1, Y4, Y1
	VPMINSD      Y1, Y5, Y1
	VPMINSD      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINSD       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint8MinAvx2Asm(x []uint8, r []uint8)
// Requires: AVX, AVX2, SSE2
TEXT ·uint8MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0xffffffffffffffff, BX
	MOVQ         BX, X0
	VPBROADCASTB X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

uint8MinBlockLoop:
	CMPQ    DX, $0x000000c0
	JL      uint8MinTailLoop
	VPMINUB (AX), Y1, Y1
	VPMINUB 32(AX), Y2, Y2
	VPMINUB 64(AX), Y3, Y3
	VPMINUB 96(AX), Y4, Y4
	VPMINUB 128(AX), Y5, Y5
	VPMINUB 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x000000c0, DX
	JMP     uint8MinBlockLoop

uint8MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint8MinDone
	VPMINUB (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000020, DX
	JMP     uint8MinTailLoop

uint8MinDone:
	VPMINUB      Y1, Y2, Y1
	VPMINUB      Y1, Y3, Y1
	VPMINUB      Y1, Y4, Y1
	VPMINUB      Y1, Y5, Y1
	VPMINUB      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINUB       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint16MinAvx2Asm(x []uint16, r []uint16)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·uint16MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0xffffffffffffffff, BX
	MOVQ         BX, X0
	VPBROADCASTW X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

uint16MinBlockLoop:
	CMPQ    DX, $0x00000060
	JL      uint16MinTailLoop
	VPMINUW (AX), Y1, Y1
	VPMINUW 32(AX), Y2, Y2
	VPMINUW 64(AX), Y3, Y3
	VPMINUW 96(AX), Y4, Y4
	VPMINUW 128(AX), Y5, Y5
	VPMINUW 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000060, DX
	JMP     uint16MinBlockLoop

uint16MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint16MinDone
	VPMINUW (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000010, DX
	JMP     uint16MinTailLoop

uint16MinDone:
	VPMINUW      Y1, Y2, Y1
	VPMINUW      Y1, Y3, Y1
	VPMINUW      Y1, Y4, Y1
	VPMINUW      Y1, Y5, Y1
	VPMINUW      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINUW       X0, X1
	MOVOU        X1, (CX)
	RET

// func uint32MinAvx2Asm(x []uint32, r []uint32)
// Requires: AVX, AVX2, SSE2, SSE4.1
TEXT ·uint32MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0xffffffffffffffff, BX
	MOVQ         BX, X0
	VPBROADCASTD X0, Y0
	VMOVDQU      Y0, Y1
	VMOVDQU      Y0, Y2
	VMOVDQU      Y0, Y3
	VMOVDQU      Y0, Y4
	VMOVDQU      Y0, Y5
	VMOVDQU      Y0, Y0

uint32MinBlockLoop:
	CMPQ    DX, $0x00000030
	JL      uint32MinTailLoop
	VPMINUD (AX), Y1, Y1
	VPMINUD 32(AX), Y2, Y2
	VPMINUD 64(AX), Y3, Y3
	VPMINUD 96(AX), Y4, Y4
	VPMINUD 128(AX), Y5, Y5
	VPMINUD 160(AX), Y0, Y0
	ADDQ    $0x000000c0, AX
	SUBQ    $0x00000030, DX
	JMP     uint32MinBlockLoop

uint32MinTailLoop:
	CMPQ    DX, $0x00000004
	JL      uint32MinDone
	VPMINUD (AX), Y1, Y1
	ADDQ    $0x00000020, AX
	SUBQ    $0x00000008, DX
	JMP     uint32MinTailLoop

uint32MinDone:
	VPMINUD      Y1, Y2, Y1
	VPMINUD      Y1, Y3, Y1
	VPMINUD      Y1, Y4, Y1
	VPMINUD      Y1, Y5, Y1
	VPMINUD      Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	PMINUD       X0, X1
	MOVOU        X1, (CX)
	RET

// func float32MinAvx2Asm(x []float32, r []float32)
// Requires: AVX, AVX2, SSE, SSE2
TEXT ·float32MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x000000007f7fffff, BX
	MOVQ         BX, X0
	VBROADCASTSS X0, Y0
	VMOVUPS      Y0, Y1
	VMOVUPS      Y0, Y2
	VMOVUPS      Y0, Y3
	VMOVUPS      Y0, Y4
	VMOVUPS      Y0, Y5
	VMOVUPS      Y0, Y0

float32MinBlockLoop:
	CMPQ   DX, $0x00000030
	JL     float32MinTailLoop
	VMINPS (AX), Y1, Y1
	VMINPS 32(AX), Y2, Y2
	VMINPS 64(AX), Y3, Y3
	VMINPS 96(AX), Y4, Y4
	VMINPS 128(AX), Y5, Y5
	VMINPS 160(AX), Y0, Y0
	ADDQ   $0x000000c0, AX
	SUBQ   $0x00000030, DX
	JMP    float32MinBlockLoop

float32MinTailLoop:
	CMPQ   DX, $0x00000004
	JL     float32MinDone
	VMINPS (AX), Y1, Y1
	ADDQ   $0x00000020, AX
	SUBQ   $0x00000008, DX
	JMP    float32MinTailLoop

float32MinDone:
	VMINPS       Y1, Y2, Y1
	VMINPS       Y1, Y3, Y1
	VMINPS       Y1, Y4, Y1
	VMINPS       Y1, Y5, Y1
	VMINPS       Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	MINPS        X0, X1
	MOVOU        X1, (CX)
	RET

// func float64MinAvx2Asm(x []float64, r []float64)
// Requires: AVX, AVX2, SSE2
TEXT ·float64MinAvx2Asm(SB), NOSPLIT, $0-48
	MOVQ         x_base+0(FP), AX
	MOVQ         r_base+24(FP), CX
	MOVQ         x_len+8(FP), DX
	MOVQ         $0x7fefffffffffffff, BX
	MOVQ         BX, X0
	VBROADCASTSD X0, Y0
	VMOVUPD      Y0, Y1
	VMOVUPD      Y0, Y2
	VMOVUPD      Y0, Y3
	VMOVUPD      Y0, Y4
	VMOVUPD      Y0, Y5
	VMOVUPD      Y0, Y0

float64MinBlockLoop:
	CMPQ   DX, $0x00000018
	JL     float64MinTailLoop
	VMINPD (AX), Y1, Y1
	VMINPD 32(AX), Y2, Y2
	VMINPD 64(AX), Y3, Y3
	VMINPD 96(AX), Y4, Y4
	VMINPD 128(AX), Y5, Y5
	VMINPD 160(AX), Y0, Y0
	ADDQ   $0x000000c0, AX
	SUBQ   $0x00000018, DX
	JMP    float64MinBlockLoop

float64MinTailLoop:
	CMPQ   DX, $0x00000004
	JL     float64MinDone
	VMINPD (AX), Y1, Y1
	ADDQ   $0x00000020, AX
	SUBQ   $0x00000004, DX
	JMP    float64MinTailLoop

float64MinDone:
	VMINPD       Y1, Y2, Y1
	VMINPD       Y1, Y3, Y1
	VMINPD       Y1, Y4, Y1
	VMINPD       Y1, Y5, Y1
	VMINPD       Y1, Y0, Y1
	VEXTRACTF128 $0x01, Y1, X0
	MINPD        X0, X1
	MOVOU        X1, (CX)
	RET
