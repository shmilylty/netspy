package core

import "fmt"

const LOGO = `
               __                        
  ____   _____/  |_  ____________ ___.__.
 /    \_/ __ \   __\/  ___/\____ <   |  |
|   |  \  ___/|  |  \___ \ |  |_> >___  |
|___|  /\___  >__| /____  >|   __// ____|
     \/     \/          \/ |__|   \/
`

func GetBanner() string {
	banner := fmt.Sprintf("%s%s\n\n", LOGO, GetVersion())
	return banner
}
