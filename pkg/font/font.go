package font

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

type Font struct {
	Family   string
	Style    string
	Filepath string
}

func HasFontFamily(fonts []string, predicate string) bool {
	return slices.Contains(fonts, predicate)
}

var fonts = make([]Font, 0)

func GetFonts() ([]Font, error) {
	if len(fonts) > 0 {
		return fonts, nil
	}

	out, err := exec.Command("fc-list").Output()
	if err != nil {
		return nil, err
	}

	text := string(out)
	lines := strings.SplitSeq(text, "\n")
	for v := range lines {
		line := strings.Split(v, ":")

		font := Font{
			Filepath: "",
			Family:   "",
			Style:    "",
		}

		if len(line) == 0 {
			continue
		}

		if len(line) > 0 {
			font.Filepath = line[0]
		}
		if len(line) > 1 {
			fam := strings.Split(line[1], ",")
			font.Family = strings.TrimSpace(fam[0])
		}
		if len(line) > 2 {
			font.Style = strings.TrimLeft(line[2], "style=")
		}

		fonts = append(fonts, font)
	}

	return fonts, nil

}

func GetFontFamilies() ([]string, error) {
	fonts, err := GetFonts()
	if err != nil {
		return nil, err
	}

	list := []string{}
	for _, v := range fonts {
		if HasFontFamily(list, v.Family) == false {
			list = append(list, v.Family)
		}
	}

	return list, nil
}

func GetFontStyles(family string) ([]string, error) {
	styles := []string{}

	fonts, err := FindFonts(family, "")
	if err != nil {
		return nil, err
	}

	for _, v := range fonts {
		if v.Family == family && v.Style != "" {
			styles = append(styles, v.Style)
		}
	}

	return styles, nil
}

func FindFonts(family string, style string) ([]Font, error) {
	fonts, err := GetFonts()

	if err != nil {
		return nil, err
	}

	list := []Font{}
	for _, v := range fonts {
		if family != "" && v.Family != family {
			continue
		}
		if style != "" && v.Style != style {
			// could be many
			if !strings.Contains(fmt.Sprintf(",%s,", style), v.Style) {
				continue
			}
		}

		// duplicated list
		list = append(list, v)
	}

	// deduplication logic
	seen := make(map[string]bool)
	result := []Font{}

	for _, item := range list {
		tag := item.Family+item.Style
		if !seen[tag] {
			seen[tag] = true
			result = append(result, item)
		}
	}

	return result, err
}
