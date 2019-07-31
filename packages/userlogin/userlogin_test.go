package userlogin

import "testing"

func TestPasswordHashing(t *testing.T) {
	original := "$2a$08$KVtTpQwdY8ZdkBcgilC2COBIIy1lOATEbmR0KmRz7hqayoG05WsTq"

	pass := passwordHashing([]byte("root@123"))

	if original != pass {
		t.Errorf("two outputs")
	}
}
