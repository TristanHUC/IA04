//go:build js
// +build js

package _map

var backupMap = `wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww                       
w                                                     wwwwww                       
w                                                     wwwwww                       
w                                                     wwwwww                       
w                                                     wwwwww                       
w                                                     wwwwww                       
w                                                     wwwwww                       
w                                                     d   dw                       
w                                                       d  w                       
w                                                     d   dw                       
w                                                          w                       
w                                                       d  w                       
w                                                     d   dw                       
w                                                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w            wwwwwwwwwwwwwwwwwwww                     wwwwww                       
w           wwwwwwwwwwwwwwwwwwwww                     wwwwww                       
w           wwwwwwwwwwwwwwwwwwwww                     wwwwww                       
w          wwwwwwwwwwwwwwwwwwwwww                     wwwwww                       
wwww    wwwwwwwwwwwwwwwwwwwwwwwwww  wwwwwwwwwwww      wwwwww              e        
w        bacccwwwwwwwwwwwwwwwwwww    w     w m w      wwwwww                       
w        bacccwwwwwwwwwwwwwwwwwww    w     w   w       wwwww                       
w        bacccwwwwwwwwwwwwwwwwwww    w     w   w        wwww     e                 
w        bacccwwwwwwwwwwwwwwwwwww    w         w         www                       
w        bacccwwwwwwwwwwwwwwwwwww              w          ww                       
w        bacctwwwwwwwwwwwwwwwwwww          w   w           w                       
w        bacctwwwwwwwwwwwwwwwwwww          w   w           w                       
w        bacctwwwwwwwwwwwwwwwwwww    wmmmmmw   w           w                       
w        bacctwwwwwwwwwwwwwwwwwww    wwwwwwwwwww           w                       
w        bacctwwwwwwwwwwwwwwwwwww    wf f f f fw           w                      e
w        bacctwwwwwwwwwwwwwwwwwww    w         w      wwwwww                       
w        bacctwwwwwwwwwwwwwwwwwww    w         w      w                   e        
w        bacctwwwwwwwwwwwwwwwwwww              w      w                            
w        bacctwwwwwwwwwwwwwwwwww               w      w                            
w        bacctwwwwwwwwwwwwwwwwww               w      w             e              
w        bacctwwwwwwwwwwwwwwwwww               w      w    w        e              
w        bacctwwwwwwwwwwwwwwwwww               w      w    w                       
w        bacctwwwwwwwwwwwwwwwwww     w         w      w    w                       
w        bacccwwwwwwwwwwwwwwwwww     w         w      w    w                       
w        bacccwwwwwwwwwwwwwwwwww     wwwwwwwwwww      w    w                       
w        bacccwwww                             w      w    w                       
w        baaaawww                              w      w    w                 e     
w        bbbbb                                 w      w    w                       
w                                              w      w    w                       
w                                              w      w    w                       
w                                              w      w    w                       
w                                              w      w    w       e               
w                                              w           w                      e
w                                              w           w                      e
w                                              w           w                      e
w                                              w           w                       
w                                              w           w                       
w                                              w           w                       
w                                              w           w                ee     
wwwwww   wwwwwwwwwwwwwwwwwwwwwwwwwwwwww   wwwwww          ww     ee            e   
w                                                        www                   e   
w                                                       wwww                       
w                                                      wwwww                       
w                                                     wwwwww                       
w                                                    wwwwwww                       
w                                                   wwwwwwww                       
w                                                  wwwwwwwww                       
w                                                 wwwwwwwwww                       
w                                                wwwwwwwwwww                       
w                                               wwwwwwwwwwww                       
w                                              wwwwwwwwwwwww                       
wwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwwww
`

func LoadMap(mapName string) (Map, error) {
	testmap := Map{}
	err := testmap.LoadFromString(backupMap)
	if err != nil {
		return testmap, err
	}
	return testmap, nil
}