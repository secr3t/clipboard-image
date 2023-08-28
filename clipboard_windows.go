//go:build windows
// +build windows

package clip_img

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func write(file string) error {
	cmd := exec.Command("powershell", "-Command",
		fmt.Sprintf(`$Image = [System.Drawing.Image]::FromFile("%s")
							$Stream = New-Object System.IO.MemoryStream
							$Image.Save($Stream, [System.Drawing.Imaging.ImageFormat]::Jpeg)
							$Clipboard = [System.Windows.Forms.Clipboard]
							$Clipboard::SetData("Preferred DropEffect", "Copy")
							$Clipboard::SetData("JPEG", $Stream.ToArray())
							$Stream.Close()`, file))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s: %s", err, string(b))
	}
	return nil
}

func read() (io.Reader, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return nil, err
	}
	f.Close()
	defer os.Remove(f.Name())

	cmd := exec.Command("PowerShell", "-Command", "Add-Type", "-AssemblyName",
		fmt.Sprintf("System.Windows.Forms;$clip=[Windows.Forms.Clipboard]::GetImage();if ($clip -ne $null) { $clip.Save('%s') };", f.Name()))
	b, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("%s: %s", err, string(b))
	}

	r := new(bytes.Buffer)
	f, err = os.Open(f.Name())
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err := io.Copy(r, f); err != nil {
		return nil, err
	}

	return r, nil
}
