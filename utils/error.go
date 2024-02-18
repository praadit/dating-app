package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
)

func SqlPanicFilter(err error, message, notfoundMsg string) (bool, error) {
	if err == nil {
		return false, nil
	}

	switch err {
	case nil:
		return false, nil
	case sql.ErrNoRows:
		if len(notfoundMsg) < 1 {
			return true, errors.New("Record not found")
		} else {
			return true, errors.New(notfoundMsg)
		}
	default:
		log.Print(err.Error())
		return false, errors.New(message)
	}
}

func FilterError(err error, message string) error {
	if err == nil {
		return nil
	}
	log.Println(fmt.Sprintf("[INTERNAL_LOG]: %s", err.Error()))
	return errors.New(message)
}

func ERR_INVALID_TOKEN() error {
	log.Println("[INTERNAL_LOG]: access token invalid")
	return errors.New("err: a001")
}

func ERR_FAILED_JWT_PARSE() error {
	log.Println("[INTERNAL_LOG]: failed to parse jwt claims")
	return errors.New("err: a002")
}

func ERR_JWT_EXPIRED() error {
	log.Println("[INTERNAL_LOG]: jwt expired")
	return errors.New("err: a003")
}

func ERR_UNKNOWN(err error) error {
	log.Println("[INTERNAL_LOG]: " + err.Error())
	return errors.New("err: uk003")
}
