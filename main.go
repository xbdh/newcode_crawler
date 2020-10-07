package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/MontFerret/ferret/pkg/compiler"
	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/cdp"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/gin-gonic/gin"
)

type NewCodeGo struct {
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Time       string    `json:"time"`
	DiscussTag string    `json:"discuss_tag"`
	Link       string    `json:"link"`
	Content    string    `json:"content"`
}

func main() {
	interviews, err := getGolangInterview()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, interview := range interviews {

		//interview.Content=strings.Replace(interview.Content,"\n","<\br>",-1)
		fmt.Println(fmt.Sprintf("%s:\n %s \n%s \n%s \n%s \n %s\n", interview.Title, interview.Author, interview.Link, interview.Time, interview.DiscussTag, strings.Replace(interview.Content,"\\n","<br>",-1)))
	}
	r:=gin.Default()
	r.LoadHTMLGlob("templates/*")
	//r.Static("/static","./static")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200,"view.gohtml",interviews)
	})
	r.Run(":8888")
}

func getGolangInterview() ([]*NewCodeGo, error) {
	query := `
				LET parturl="http://www.nowcoder.com/discuss/tag/640?type=2&order=3&pageSize=30&expTag=0&query=&page="
				//LET index =['1','2','3','4','5','6']
				LET index =['21','22']
				
				let allarticle=(

					for i in index
					
					    LET doc = DOCUMENT(parturl+i, {
					        driver: "cdp"
					        })
					
					    LET articles = ELEMENTS(doc, '.discuss-detail')
					    LET links = (
					        FOR article IN articles
					            let urlmap=element(article,'a')
					            let url= "http://newcoder.com" + urlmap.attributes.href
					        
					            let  newurl= left(url,34) 
					            RETURN newurl
					    )
					
					    LET inner =(
					        FOR link IN links
					    
					            NAVIGATE(doc, link, 20000)
					            click(doc,'.pop-close')
					             WAIT_ELEMENT(doc, '.nk-content', 5000)
					    
					            LET texter = ELEMENT(doc, '.nk-content')
					            LET title = ELEMENT(texter, '.post-title')
					            LET name  = ELEMENT(texter, '.post-name')
					            LET time = ELEMENT(texter, '.post-time')
					            LET content =ELEMENT(texter, '.nc-post-content')
					            LET discusstag=ELEMENT(texter, '.discuss-tags-mod')
					            
					            RETURN   {
					                title: title.innerText,
					                author : name.innerText,
					                time : time.innerText,
					                discuss_tag: discusstag.innerText,
					                link : link,
					                content : content.innerText
					               
					            }
					    )
					
					    RETURN inner
				)
				
				let interview=(
				        
				    for out in allarticle
				        for each in out
				            return each
				)
				return interview
		`

	comp := compiler.New()

	program, err := comp.Compile(query)

	if err != nil {
		return nil, err
	}

	// create a root context
	ctx := context.Background()

	// enable HTML drivers
	// by default, Ferret Runtime does not know about any HTML drivers
	// all HTML manipulations are done via functions from standard library
	// that assume that at least one driver is available
	ctx = drivers.WithContext(ctx, cdp.NewDriver(),drivers.AsDefault())
	ctx = drivers.WithContext(ctx, http.NewDriver())

	out, err := program.Run(ctx)
	fmt.Printf("%s\n",out)
	if err != nil {
		return nil, err
	}
	fmt.Println()

	res := make([]*NewCodeGo ,0,50)
	//res=repalcebr(res)
	err = json.Unmarshal(out, &res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

//func repalcebr(res []*NewCodeGo)  []*NewCodeGo  {
//
//	for _, n := range res{
//		    n.Content=  strings.Replace(n.Content,"\\n","<br>",-1)
//	}
//	return res
//}