package executer

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Execute(id string, code string, input string) (string, error){
	dir := filepath.Join("functions", id)
	os.Mkdir(dir, 0755)
	file := filepath.Join(dir, "main.go")

	err := os.WriteFile(file, []byte(code), 0644)

	if err != nil{
		return "", err
	}
	cmd := exec.Command("go", "run", file, input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}