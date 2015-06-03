package zutil

import (
    "fmt"
)

func GetNumericVersion() (string) {
    return fmt.Sprintf("%d.%d.%d", ZAGT_MAJOR, ZAGT_MINOR, ZAGT_BUILD)
}

