package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Instructions []byte

type Opcode byte

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTruthy
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
)

type Definition struct {
	name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
}

// Lookup 传入opcode的byte
// 得到opcode的定义
func Lookup(op byte) (*Definition, error) {
	res, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return res, nil
}

// Make 传入操作码、操作数，得到字节码,其中操作数使用大端编码
//
//	 for example :
//		Make(OpConstant, []int{65534})
//	 返回[]byte{byte(OpConstant), 255, 254}
func Make(op Opcode, operand ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	var instruction []byte
	if len(def.OperandWidths) == 0 {
		instruction = make([]byte, 1)
		instruction[0] = byte(op)
	} else {
		instructionLen := def.OperandWidths[0] + 1
		instruction = make([]byte, instructionLen)
		instruction[0] = byte(op)
	}
	offset := 1
	for i, o := range operand {
		with := def.OperandWidths[i]
		switch with {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += with
	}
	return instruction
}

// ReadOperands : Make 的逆过程
// 传入opcode的定义、字节码指令
// 返回字节码的操作数operands和指令长度
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUnit16(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUnit16 辅助函数
// 将[]byte转化为uint16
func ReadUnit16(ins []byte) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func (ins Instructions) String() string {
	var out bytes.Buffer
	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "Error:%s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))
		i += 1 + read
	}
	return out.String()
}

// fmtInstruction 辅助函数，用于将指令instructions格式化输出
func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandsCount := len(def.OperandWidths)
	if len(operands) != operandsCount {
		return fmt.Sprintf("ERROR:operands len %d does not match defined %d\n", len(operands), operandsCount)
	}

	switch operandsCount {
	case 0:
		return def.name
	case 1:
		return fmt.Sprintf("%s %d", def.name, operands[0])
	}

	return fmt.Sprintf("ERROR:unhandled operandCount for %s\n", def.name)
}
