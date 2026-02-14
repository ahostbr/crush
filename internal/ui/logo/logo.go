// Package logo renders the Kuroryuu brand art and wordmark.
package logo

import (
	"fmt"
	"image/color"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/ahostbr/crush/internal/ui/styles"
	"github.com/charmbracelet/x/ansi"
)

const diag = `╱`

// Opts are the options for rendering the Kuroryuu title art.
type Opts struct {
	FieldColor   color.Color // diagonal lines
	TitleColorA  color.Color // left gradient ramp point
	TitleColorB  color.Color // right gradient ramp point
	CharmColor   color.Color // Kuroryuu text color
	VersionColor color.Color // Version text color
	Width        int         // width of the rendered logo, used for truncation
}

// dragonArt is the Kuroryuu dragon ASCII art.
const dragonArt = `                                            -=++++++++=+=
                                    ====    -==+++++++======--==
                                #+=+                     ====+++++++
                             ### %%                           ++******#%%
                          ++*     %%        **                   ****#%%%%%
                        ++        %%%        ##                     #%%%%%%%%%
                     +             #%%%%     %%##                  *** %%%%%%%%%
                                     %%%      %%%                  +  #  #%##%%%%%
                                     %%%%      %%%                ++  #    %##%%%%%%
                                      %%%%     %%%                  =##     #***##*#%
                                *#      %%%% %% %%%  %             **         *++++*%%  * **
                                  %%%%  #%%%%%%% %%%%%%           #            ++*##%***     ++
                                   ##%%####*%%%%%% %%%%%         #              *########* **    *
                      ++=++     ##%  ##%%##%%%%%%%%%%%%%% #      %              **####**#+=+#*+
                      ++=         ##%%%%%%##%%%%%###%%%%%%###   *%             =+**####*****##=+
                        =       %%%%%%%%%%###%%%%%%%%%%#######   #               +**########+==+
                       ***##***####%%%%%%%%%###%%%%%%%##******#* **             ++**###%####*++#
                     ****%%%%*+*#%%%%%%%%%%#####%%%%%###******#   ##             **+**#%%%%%##*#
                     **###%%%%%##%%#%%%%%%%#%##%%%%%%%###***##%%%%%##          ++**--*##%%###***
                       *#%%%%%%%#%%%%%%%%%%%%%%%%%%###%#****##%%%%%%%#         ++++--+*#%%%%%#**
               **++   +**#*#%%%%%%%%%%%%%##########%#####+*    #%%%%#*       == +*+==*###%%%%%****#
                +++**++***+*%%%%%%%#%%%#####*#   %%%  *#%%**     ###*       === +*****###%%%%%**#**
                  ###***   =*#%%%%%########             %%%#%##*           === +***#########%%#####
               ######*+=     +#%%%%%##***+ %        ***# #%%#%%%*           ++++***##%%%###*#%%#*##
               #**%%%*++   --=+#%%%%%###*#               %%%%%%%      +     +++**#####%%%%##%%%####
               ## #** **+  ===+*######*==           %##%%#%%%%%   *   **  ++*#####%%##%%###%%%%###
                  %%% +**   ==+*####*=======       #             **+*+= ****##**#%%#######%#%%%###
      +             #**++   =+++****+#+===+****                 =+****#*****#########%%%%%%%%%%##
                     ==     =++==+*++==**####++                 ==++**#####%##%#*#########%%%%#*
       +               ++    =+*++++==*#%%%%###%              ===+++++**%%#%%%%%%%###%###%#%##*
                             +*##*+*######%##                  %%#*=+#%%%##%%%%%#####%####%%#***
                              ###**#####%%%%%##+      #%*  %%%%%%%%#%%%%%%%%%%%%%%##+*#%%%%%#*
         ##                   ###%%%%%%%%%%%%%%*+=%%%%%%%#%%%%##%%%%%%%%%%%%%%%%%%%%%***##%##*
   *  #                         #%#%%%%%%%%%%%%%##%%%#%%%%####**%%%%%%%%%%%%%%%%%%%%%%####%%#
 *++ ++  #      +                ###*##%%%%%%%%##***#%%%%#*+*###%%%%%#%%%%%%%%%%%%%%%%%%%#%%*=
   *++  #   #    +             =+++****##%%%%%##**++*#%%%#*++%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#*===
     ==   * ==#%% #            -=   *##%%%%%####***++++*+=*###%%%%%%####%%%%%%%%%%%%%%%%%%#*
      ==+****+ %%###%%            **++*%%%%%##*####*++++*#####%#%%%%%%%%%%%%%%%%%%%%%%%%%##*
       **** ++  %%%%%%%              *+  %#%%######**+++*####%%%%%%%%%#*#%%%%%%%%%%####%###
             **%%%%%%%%%%#*                  **#****+++*##*##%%%%%%%%#++**#%%%%%%%%##*++
                #%%%%%%%%%%  =                    ==   #   %%%%% #%###===*#%%%%%%%%#*+=
                +%%%%%%%%%%#                        %%%%%%%%%%%  %%  *=-*#%%%%%%%%%*+=
                  %%%%%%%%#                            %%%%%%%%        *#%%%%%%%%%*#*=
                   %%%%%%%%#                          %%%%%% %%       #%%%%%%%%%#*++
                      ##%%#*              *           %%  %%  %#    %%%%%%%%%%#**==
                          ##**##                           %#     ##%%%##%%%%##*=
                            ####### *                         ###%%%%%%%%%%++*+*
                              ######%#*%%         #   ***==*%%%*+*#*##%%%%
                                 ####%%%%%%##%#%%%%%%%%###*#%%%*     +#
                                   +**#%%%%%%%%%%%%%%%%%#%%%%%%#==
                                           %%%%%%%%%%%%%%%%#%`

// DragonArt returns the dragon ASCII art for display on landing pages.
func DragonArt() string {
	return dragonArt
}

// Render renders the Kuroryuu header logo.
func Render(s *styles.Styles, version string, compact bool, o Opts) string {
	fg := func(c color.Color, str string) string {
		return lipgloss.NewStyle().Foreground(c).Render(str)
	}

	title := "黒竜 KURORYUU"
	titleWidth := lipgloss.Width(title)

	// Apply gradient to title
	gradTitle := styles.ApplyBoldForegroundGrad(s, title, o.TitleColorA, o.TitleColorB)

	// Version
	versionStr := fg(o.VersionColor, version)

	if compact {
		// Compact sidebar version
		field := fg(o.FieldColor, strings.Repeat(diag, titleWidth))
		return strings.Join([]string{field, gradTitle, field, ""}, "\n")
	}

	// Wide version
	gap := max(1, titleWidth-lipgloss.Width(version))
	metaRow := fg(o.CharmColor, " Kuroryuu") + strings.Repeat(" ", gap) + versionStr
	metaWidth := lipgloss.Width(metaRow)
	contentWidth := max(titleWidth, metaWidth)

	// Left field
	const leftWidth = 4
	leftFieldRow := fg(o.FieldColor, strings.Repeat(diag, leftWidth))
	fieldHeight := 3 // meta + title + bottom line
	leftField := new(strings.Builder)
	for range fieldHeight {
		fmt.Fprintln(leftField, leftFieldRow)
	}

	// Right field
	rightWidth := max(10, o.Width-contentWidth-leftWidth-2)
	rightField := new(strings.Builder)
	for i := range fieldHeight {
		w := rightWidth - i
		if w < 0 {
			w = 0
		}
		fmt.Fprintln(rightField, fg(o.FieldColor, strings.Repeat(diag, w)))
	}

	// Assemble center content
	bottomLine := fg(o.FieldColor, strings.Repeat(diag, contentWidth))
	center := strings.Join([]string{metaRow, gradTitle, bottomLine}, "\n")

	const hGap = " "
	logo := lipgloss.JoinHorizontal(lipgloss.Top, leftField.String(), hGap, center, hGap, rightField.String())
	if o.Width > 0 {
		lines := strings.Split(logo, "\n")
		for i, line := range lines {
			lines[i] = ansi.Truncate(line, o.Width, "")
		}
		logo = strings.Join(lines, "\n")
	}
	return logo
}

// SmallRender renders a compact version of the Kuroryuu logo.
func SmallRender(t *styles.Styles, width int) string {
	title := t.Base.Foreground(t.Secondary).Render("黒竜")
	title = fmt.Sprintf("%s %s", title, styles.ApplyBoldForegroundGrad(t, "Kuroryuu", t.Secondary, t.Primary))
	remainingWidth := width - lipgloss.Width(title) - 1
	if remainingWidth > 0 {
		lines := strings.Repeat("╱", remainingWidth)
		title = fmt.Sprintf("%s %s", title, t.Base.Foreground(t.Primary).Render(lines))
	}
	return title
}
