package hash

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/alexedwards/argon2id"
)

type hashProvider struct {
}

var params_ = argon2id.Params{
	Memory:      64 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  32,
	KeyLength:   32,
}

var emptyPassword_ string

func init() {
	emptyPassword_ = HashProvider().HashPassword("")
}

func HashProvider() hashProvider {
	return hashProvider{}
}

func (h hashProvider) log(format string, v ...interface{}) {
	log.Printf("|HASH> "+format, v...)
}

func (h hashProvider) EmptyPassword() string {
	return emptyPassword_
}

func (h hashProvider) CalibrateMemoryParam(maxExecutionTime time.Duration) (err error) {
	h.log("Calibrating memory params for password hash")
	h.log("Current limit: %s", maxExecutionTime.String())

	var hashTime time.Duration
	var hash string
	par := params_
	par.Memory = 1024
	for err == nil {
		par.Memory *= 2
		h.log("Trying %s", formatBytes(int64(par.Memory)))
		hashTime, hash, err = hashingDuration(
			&par,
			maxExecutionTime+50*time.Millisecond,
		)

		if err == nil {
			emptyPassword_ = hash
			h.log("Passed %s in %s", formatBytes(int64(par.Memory)), hashTime.String())
		}
	}
	h.log("Failed %s after %s", formatBytes(int64(par.Memory)), hashTime.String())

	params_.Memory = par.Memory / 2
	h.log("Done calibrating memory params for password hash")
	h.log("Using params: %+v", params_)
	return nil
}

func (h hashProvider) HashPassword(password string) string {
	hash, _ := argon2id.CreateHash(password, &params_)

	return hash
}

func (h hashProvider) VerifyPassword(password, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)

	if err != nil {
		return false
	}

	return match
}

func hashingDuration(params *argon2id.Params, timeout time.Duration) (time.Duration, string, error) {
	type result struct {
		success bool
		hash    string
	}

	processDone := make(chan result)
	start := time.Now()
	go func() {
		hash, err := argon2id.CreateHash(
			"",
			params,
		)
		processDone <- result{
			success: err == nil,
			hash:    hash,
		}
	}()

	select {
	case <-time.After(timeout):
		return time.Since(start), "", errors.New("timeout")
	case res := <-processDone:
		if !res.success {
			return time.Since(start), "", errors.New("something went wrong")
		}
		return time.Since(start), res.hash, nil
	}
}

func formatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf(
		"%.1f %ciB",
		float64(b)/float64(div),
		"KMGTPE"[exp],
	)
}
