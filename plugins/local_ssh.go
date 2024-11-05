package plugins

import (
	"bufio"
	"bytes"
	"context"
	"github.com/chris-cmsoft/conftojson/pkg"
	"github.com/open-policy-agent/opa/rego"
	"os/exec"
)

func NewLocalSSH() *LocalSSH {
	return &LocalSSH{}
}

type LocalSSH struct {
	data map[string]interface{}
}

func (l *LocalSSH) PrepareForEval(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "ssh", "root@jgc", "sshd", "-T")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(stdout)
	scanner := bufio.NewScanner(buf)

	sshConfigMap, err := pkg.ConvertConfToMap(scanner)
	if err != nil {
		return err
	}

	l.data = sshConfigMap
	return nil
}

func (l *LocalSSH) Evaluate(ctx context.Context, query rego.PreparedEvalQuery) (rego.ResultSet, error) {
	result, err := query.Eval(ctx, rego.EvalInput(l.data))
	return result, err
}
