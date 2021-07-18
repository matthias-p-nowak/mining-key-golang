package main

import (
  "fmt"
  "flag"
  "os"
  "crypto/sha256"
  "math/bits"
  "time"
)

var (
  salt=make(chan string)
  alphabet="0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func getSalt(first string, depth, level int){
  level=level+1
  l:=len(alphabet)
  for i:=0;i<l;i++{
    str:=first+string(alphabet[i])
    if level >= depth {
      salt <- str
    } else {
      getSalt(str,depth,level);
    }
  }
}

func getAllSalt(){
  for i:=1;i<=32;i++{
    getSalt("",i,0);
  }
  close(salt);
}

func main(){
  fmt.Println("starting")
  defer fmt.Println("all done")
  lz:=flag.Int("b", 24, "number of leading zeros");
  
  flag.Parse()
  go getAllSalt()
  args:=flag.Args()
  if len(args)<1 {
    fmt.Printf("have no key to mine")
    os.Exit(2)
  }
  key:=args[0]
  bestkey:=""
  for str:= range salt {
    str=key+"-"+str
    bb:=sha256.Sum256([]byte(str))
    for k:=0;k<len(bb);k++{
      if bb[k]!=0 {
        k=k*8+bits.LeadingZeros8(bb[k]);
        if k>= *lz {
          bestkey=str
          t:=time.Now()
          fmt.Printf("%s: %s -> %d\n",t.Format("15:04:05"),bestkey,k);
        }
        break
      }
    }
  }
}
