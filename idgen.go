package popup

//id generator

import (
	"errors"
	"fmt"
	// "log"
	"sync"
	"time"
	"github.com/sumory/baseN4go"
)

const (
	twepoch            = int64(1288834974657)
	workerIdBits       = uint(5)
	datacenterIdBits   = uint(5)
	maxWorkerId        = -1 ^ (-1 << workerIdBits)     //31
	maxDatacenterId    = -1 ^ (-1 << datacenterIdBits) //31
	sequenceBits       = uint(12)
	workerIdShift      = sequenceBits                                   //12
	datacenterIdShift  = sequenceBits + workerIdBits                    //17
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits //22
	sequenceMask       = -1 ^ (-1 << sequenceBits)                      //4095
)

type IdWorker struct {
	sequence      int64
	lastTimestamp int64
	workerId      int64
	datacenterId  int64
	mutex *sync.Mutex
	baseN *baseN4go.BaseN
}

func NewIdWorker(workerId, datacenterId int64) (*IdWorker, error) {
	idWorker := &IdWorker{}
	if workerId > maxWorkerId || workerId < 0 {
		return nil, errors.New(fmt.Sprintf("worker Id: %d error", workerId))
	}
	if datacenterId > maxDatacenterId || datacenterId < 0 {
		return nil, errors.New(fmt.Sprintf("datacenter Id: %d error", datacenterId))
	}
	idWorker.workerId = workerId
	idWorker.datacenterId = datacenterId
	idWorker.lastTimestamp = -1
	idWorker.sequence = 0
	idWorker.mutex = &sync.Mutex{}
	err,baseN := baseN4go.NewBaseN(int8(62))
	if err!=nil{
		return nil, errors.New(fmt.Sprintf("baseN4go error"))
	}
	idWorker.baseN = baseN
	return idWorker, nil
}

func timeGen() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

func (id *IdWorker) NextId() (int64, error) {
	id.mutex.Lock()
	defer id.mutex.Unlock()
	timestamp := timeGen()
	if timestamp < id.lastTimestamp {
		return 0, errors.New(fmt.Sprintf("Clock error! Wait for %d milliseconds", id.lastTimestamp-timestamp))
	}
	if id.lastTimestamp == timestamp {
		id.sequence = (id.sequence+1)&sequenceMask
		if id.sequence == 0 {
			timestamp = tilNextMillis(id.lastTimestamp)
		}
	} else {
		id.sequence = 0
	}
	id.lastTimestamp = timestamp
	return ((timestamp-twepoch)<<timestampLeftShift)|(id.datacenterId<<datacenterIdShift)|(id.workerId<<workerIdShift)|id.sequence, nil
}

func (id *IdWorker) Encode(source int64) string {
	_,result := id.baseN.Encode(source)
	return result
}

func (id *IdWorker) Decode(source string) int64 {
	_,result := id.baseN.Decode(source)
	return result
}

func (id *IdWorker) ShortId() string{
	newId,_:=id.NextId()
	return id.Encode(newId)
}

//func main() {
//	idWorker, err := NewIdWorker(1, 1)
//	if err != nil {
//		fmt.Println("Fatal error")
//	} else {
//		for i := 0; i < 10; i++ {
//			newId, _ := idWorker.NextId()
//			fmt.Println(newId, idWorker.Encode(newId), idWorker.Decode(idWorker.Encode(newId)))
//		}
//	}
//}
