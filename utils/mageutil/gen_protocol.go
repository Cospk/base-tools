package mageutil

import (
	"archive/zip"
	"bufio"
	"fmt"
	"github.com/magefile/mage/sh"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// 确保工具已安装
func ensureToolsInstalled() error {
	tools := map[string]string{
		"protoc-gen-go": "https://github.com/golang/protobuf/tree/master/protoc-gen-go@latest",
	}

	// 根据操作系统设置 GOBIN，Windows 需要不同的默认路径
	var targetDir string
	if runtime.GOOS == "windows" {
		targetDir = filepath.Join(os.Getenv("USERPROFILE"), "go", "bin")
	} else {
		targetDir = "/usr/local/bin"
	}

	os.Setenv("GOBIN", targetDir)

	for tool, path := range tools {
		if _, err := exec.LookPath(filepath.Join(targetDir, tool)); err != nil {
			fmt.Printf("安装 %s 到 %s...\n", tool, targetDir)
			if err := sh.Run("go", "install", path); err != nil {
				return fmt.Errorf("安装 %s 失败: %s", tool, err)
			}
		} else {
			fmt.Printf("%s 已安装在 %s。\n", tool, targetDir)
		}
	}

	if _, err := exec.LookPath(filepath.Join(targetDir, "protoc")); err == nil {
		fmt.Println("protoc 已安装。")
		return nil
	}

	fmt.Println("安装 protoc...")
	return installProtoc(targetDir)
}

// 安装 protoc
func installProtoc(installDir string) error {
	version := "26.1"
	baseURL := "https://github.com/protocolbuffers/protobuf/releases/download/v" + version
	archMap := map[string]string{
		"amd64": "x86_64",
		"386":   "x86",
		"arm64": "aarch64",
	}
	protocFile := "protoc-%s-%s.zip"

	osArch := runtime.GOOS + "-" + getProtocArch(archMap, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		osArch = "win64" // 假设 64 位，32 位使用 "win32"
	}
	fileName := fmt.Sprintf(protocFile, version, osArch)
	url := baseURL + "/" + fileName

	fmt.Println("URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "protoc-*.zip")
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("tmp ", tmpFile.Name(), "install  ", installDir)
	return unzip(tmpFile.Name(), installDir)
}

// 解压缩文件
func unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// 获取 protoc 架构
func getProtocArch(archMap map[string]string, goArch string) string {
	if arch, ok := archMap[goArch]; ok {
		return arch
	}
	return goArch
}

// 生成协议文件
func Protocol() error {
	if err := ensureToolsInstalled(); err != nil {
		fmt.Println("错误 ", err.Error())
		os.Exit(1)
	}

	moduleName, err := getModuleNameFromGoMod()
	if err != nil {
		fmt.Println("错误获取模块名从 go.mod: ", err.Error())
		os.Exit(1)
	}

	protoPath := "./pkg/protocol"
	dirs, err := os.ReadDir(protoPath)
	if err != nil {
		fmt.Println("错误 ", err.Error())
		os.Exit(1)
	}

	for _, dir := range dirs {
		if dir.IsDir() {
			if err := compileProtoFiles(protoPath, dir.Name(), moduleName); err != nil {
				fmt.Println("错误 ", err.Error())
				os.Exit(1)
			}
		}
	}
	return nil
}

// 编译协议文件
func compileProtoFiles(basePath, dirName, moduleName string) error {
	protoFile := filepath.Join(basePath, dirName, dirName+".proto")
	outputDir := filepath.Join(basePath, dirName)
	module := moduleName + "/pkg/protocol/" + dirName
	args := []string{
		"--go_out=plugins=grpc:" + outputDir,
		"--go_opt=module=" + module,
		protoFile,
	}
	fmt.Printf("编译 %s...\n", protoFile)
	if err := sh.Run("protoc", args...); err != nil {
		return fmt.Errorf("编译 %s 失败: %s", protoFile, err)
	}
	return fixOmitemptyInDirectory(outputDir)
}

// 修复 omitempty
func fixOmitemptyInDirectory(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.pb.go"))
	if err != nil {
		return fmt.Errorf("获取 .pb.go 文件列表失败: %s", err)
	}
	fmt.Printf("修复 omitempty 在目录 %s...\n", dir)
	for _, file := range files {
		fmt.Printf("修复 omitempty 在文件 %s...\n", file)
		if err := RemoveOmitemptyFromFile(file); err != nil {
			return fmt.Errorf("替换 omitempty 失败: %s", err)
		}
	}
	return nil
}

// 从文件中移除 omitempty
func RemoveOmitemptyFromFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("打开文件失败: %s", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, ",omitempty", "")
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取文件失败: %s", err)
	}

	return writeLines(lines, filePath)
}

// 将多行内容写入指定文件
func writeLines(lines []string, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("创建文件失败: %s", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range lines {
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("写入文件失败: %s", err)
		}
	}
	return w.Flush()
}

// 从 go.mod 文件中获取模块名
func getModuleNameFromGoMod() (string, error) {
	file, err := os.Open("go.mod")
	if err != nil {
		return "", fmt.Errorf("打开 go.mod 文件失败: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "module ") {
			// 假设该行形如："module github.com/user/repo"
			return strings.TrimSpace(strings.TrimPrefix(line, "module")), nil
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("读取 go.mod 文件失败: %v", err)
	}

	return "", fmt.Errorf("module directive not found in go.mod")
}
