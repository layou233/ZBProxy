package main

import "log"

/*
This code reference from wiki.vg.
You can browse more information about Minecraft protocol there.
VarInt and VarLong: https://wiki.vg/Protocol#Definitions
Copyright (C) 2021 layou233.
*/

//const maxVarintBytes = 10 // maximum length of a varint

// EncodeVarint returns the byte slice of the encoded Varint value.
func EncodeVarint(x int) []byte {
	result := make([]byte, 5)
	_tk := true // Simulate do-while
	for numWrite := 0; (x != 0) || _tk; numWrite++ {
		_tk = false
		temp := (byte)(x & 0b01111111)
		x = (int)((uint32)(x) >> 7) // unsigned right shift assignment
		if x != 0 {
			temp |= 0b10000000
		}
		result[numWrite] = temp
	}
	var length int8
	for length = 4; result[length] == 0; length-- {
	}
	return result[:length+1]
}

// DecodeVarint decodes the Varint and returns a int typed value.
// This will read bytes from the given index and read most at 5 bytes.
func DecodeVarint(buf []byte, index int) (result int, numRead int) {
	numRead, result = 0, 0
	var read byte
	_tk := true // Simulate do-while
	for ((read & 0b10000000) != 0) || _tk {
		_tk = false
		read = buf[index+numRead]
		value := (int)(read & 0b01111111)
		result |= value << (7 * numRead)
		numRead++
		if numRead > 5 {
			log.Println("Found a wrong Varint that bigger than 5. What was happened?")
			break
		}
	}
	return
}

/*
func unsignedRightShift(x, rightShiftAmount int) int {
	return (int)((uint32)(x) >> rightShiftAmount)
}*/
