package main

import (
  "fmt"
  "flag"
  "os"
  "crypto/sha256"
  "math/bits"
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
  
  flag.Parse()
  go getAllSalt()
  args:=flag.Args()
  if len(args)<1 {
    fmt.Printf("have no key to mine")
    os.Exit(2)
  }
  key:=args[0]
  bestkey:=""
  blz:=0
  for str:= range salt {
    str=key+"-"+str
    bb:=sha256.Sum256([]byte(str))
    for k:=0;k<len(bb);k++{
      if bb[k]!=0 {
        k=k*8+bits.LeadingZeros8(bb[k]);
        if k>=blz {
          bestkey=str
          blz=k
          fmt.Printf("%s -> %d\n",bestkey,blz);
        }
        break
      }
    }
  }
}
