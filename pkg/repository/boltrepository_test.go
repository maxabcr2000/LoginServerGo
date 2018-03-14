package repository_test

import (
	"testing"
	. "github.com/maxabcr2000/LoginServerGo/pkg/repository"
)

const(
	TestDB = "Test.db"
	TestBucketName = "TestBucket"
	TestKey = "TestKey"
	TestValue = "TestValue"
)

func TestRepository(t *testing.T) {
	repo,err:= CreateBoltRepository(TestDB);
	if err!=nil{
		t.Error(err)
	}

	defer repo.Close()

	err=repo.CreateBucket(TestBucketName)
	if err!=nil{
		t.Error("Error occurred while calling CreateBucket(): ", err)
	}

	err=repo.SaveMessage(TestBucketName, TestKey, TestValue)
	if err!=nil{
		t.Error("Error occurred while calling SaveMessage(): ", err)
	}

	value,err:=repo.ReadMessage(TestBucketName, TestKey)
	if err!=nil{
		t.Error("Error occurred while calling ReadMessage(): ", err)
	}

	if value!=TestValue{
		t.Errorf("Expected to get value: %s but get %s instead.", TestValue, value)
	}
}