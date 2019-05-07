package main
//script for monitoring webservices in the local network
// input: sites.txt
// ouput: log.Text


import "fmt"
import "os"
import "net/http"
import "time"
import "github.com/gookit/color"
import "io"
import "bufio"
import "strings"
import "runtime"



const CHECKS = 10
const TIME_INTERVAL = 15


func splashScreen(){
// not really needed right now, but maybe other features can be added later.

  fmt.Println("Options:")
  fmt.Println("1 - Start")
  fmt.Println("0 - Exit")
}


func readInput() int {
  var entry int
  fmt.Print("input: ")
  fmt.Scanf("%d",&entry)
  fmt.Println("value: ",entry)
  return entry
}


func main() {

  splashScreen()

  for{
    entry:=readInput()
    if  entry!=1 && entry!=0 {
      fmt.Println("Invalid Input! ")
      os.Exit(-1)
    }else{
      if entry==1{
        fmt.Println("Monitor mode")
        initMonitor()
      }else if entry==0{
        fmt.Println("Exit!")
        os.Exit(0)
      }
    }
  }
}

//c chan string
func testConnectivity(c chan string){

  fileSlice := strings.Split(<-c, ",")
  for j:=0; j<CHECKS; j++{
        for _,site := range fileSlice{
          if site!=""{
              resp,err:=http.Get(site)
              if err!=nil{
                  color.Error.Println(err)
                  writeLog(site,false,111)
                }else{
                  writeLog(site, resp.StatusCode==200 ,resp.StatusCode)
                  if resp.StatusCode == 200{
                    color.Info.Println("Site: "+site+" status: OK.")
                  } else {
                    color.Warn.Println("Site: "+site+" not reachable. Status Code: ",resp.StatusCode," ➤ Status:",resp.Status)
                  }
                }
          }
        }
      time.Sleep(TIME_INTERVAL*time.Second)
  }

}



func initMonitor(){
    NUMBER_CPUS:=runtime.NumCPU()
    sites:=readFile()
    step:=len(sites)/NUMBER_CPUS

    for i:=0; i<len(sites);i=i+step {

      reg:= strings.Join(sites[i:i+step], ",") // passing message as string
      message:=make(chan string)
      go testConnectivity(message)

      message <-reg
    }
}

func readFile() [] string{

  var sites []string
  file,err:=os.Open("sites.txt")
   if err!=nil{
     color.Error.Println(err)
   }

  reader := bufio.NewReader(file)
  for{
    line,err := reader.ReadString('\n') // single quote for byte
    line = strings.TrimSpace(line)
    if err==io.EOF{
     break
    }
    sites =append(sites,line)
  }
  file.Close()

  return sites
}


func writeLog(site string, status bool, code int){

  file,err := os.OpenFile("log.txt",os.O_RDWR|os.O_CREATE|os.O_APPEND,0666)
  if err!=nil{
    color.Error.Println(err)
  }
  t :=time.Now()
  s_t:=t.Format("2006-01-02,15:04:05")
  if status{
    file.WriteString(s_t+","+site + ",UP\n")
  } else {
    file.WriteString(fmt.Sprintf("%s,%s,DOWN,%d\n",s_t,site,code))
  }
  file.Close()
}
