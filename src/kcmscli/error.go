/*
 * k8s-cms
 * kcmscli - k8s-cms comand line client
 * Error Handling
 */
package main

import (
	"fmt"
	"os"
)

/* Error Handling  */
// kill program due to a fatal error detailed by given messagea
func die(message string) {
	fmt.Printf("FATAL: %s\n", message)
	os.Exit(1)
}
