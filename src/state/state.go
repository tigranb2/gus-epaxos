package state

import (
	"sync"
	//"fmt"
	//"code.google.com/p/leveldb-go/leveldb"
	//"encoding/binary"
)

type Operation uint8

const (
	NONE Operation = iota
	PUT
	GET
	DELETE
	RLOCK
	RMW
)

type Value int64

const NIL Value = 0

type Key int64

type Command struct {
	Op Operation
	K  Key
	V  Value
}

type State struct {
	mutex *sync.Mutex
	Store map[Key]Value
}

func InitState() *State {
	/*
	   d, err := leveldb.Open("/Users/iulian/git/epaxos-batching/dpaxos/bin/db", nil)

	   if err != nil {
	       fmt.Printf("Leveldb open failed: %v\n", err)
	   }

	   return &State{d}
	*/

	return &State{new(sync.Mutex), make(map[Key]Value)}
}

func Conflict(gamma *Command, delta *Command) bool {
	if gamma.K == delta.K {
		if gamma.Op == PUT || gamma.Op == RMW || delta.Op == PUT || delta.Op == RMW {
			return true
		}
	}
	return false
}

func ConflictBatch(batch1 []Command, batch2 []Command) bool {
	for i := 0; i < len(batch1); i++ {
		for j := 0; j < len(batch2); j++ {
			if Conflict(&batch1[i], &batch2[j]) {
				return true
			}
		}
	}
	return false
}

func IsRead(command *Command) bool {
	return command.Op == GET
}

func (c *Command) Execute(st *State) Value {
	//fmt.Printf("Executing (%d, %d)\n", c.K, c.V)

	//var key, value [8]byte

	//    st.mutex.Lock()
	//    defer st.mutex.Unlock()

	switch c.Op {
	case PUT:
		/*
		   binary.LittleEndian.PutUint64(key[:], uint64(c.K))
		   binary.LittleEndian.PutUint64(value[:], uint64(c.V))
		   st.DB.Set(key[:], value[:], nil)
		*/

		st.Store[c.K] = c.V
		return c.V

	case GET:
		if val, present := st.Store[c.K]; present {
			return val
		}
	case RMW:
		if val, present := st.Store[c.K]; present {
			val += 1 // modify
			st.Store[c.K] = val
			return val
		} else {
			val = 0  // default value read
			val += 1 // modify
			st.Store[c.K] = val
			return val
		}
	}

	return NIL
}
