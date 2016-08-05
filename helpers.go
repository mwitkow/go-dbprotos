// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package dbprotos

import "fmt"

var (
	UnknownFieldCallback = func(messageName string,  propertyName string) {
		fmt.Printf(`YOLO callback on %v and prop %s`, messageName, propertyName)
	}
)
