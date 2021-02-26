package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"
)

var Enctmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<body>
		<form action="/Enc">
			<label for="num"> 暗号化します >> </input>
			<input type="number" name="num" min="3" value="1">
			<input type="submit" value="決定">
		</form>
		<h1>暗号結果</h1>
		{{.}}
	</body>
</html>`))

var Dectmpl = template.Must(template.New("index").Parse(`<!DOCTYPE html>
<html>
	<body>
		<form action="/Dec">
			<label for="num2"> 復号化します  </input>
<br>暗号文   >>
			<input type="number" name="num" value="1">
<br>公開鍵 N >>
			<input type="number" name="N" value="1">
<br>公開鍵 e >>
			<input type="number" name="e"  value="1">
<br>秘密鍵 d >>
			<input type="number" name="d" value="1">
			<input type="submit" value="決定">
		</form>
		<h1>復号結果</h1>
		{{.}}
	</body>
</html>`))

func isprime1(p int) bool {
	//pが1,偶数だったらfalseを返す
	switch {
	case p == 1 || p%2 == 0:
		return false
	case p == 2:
		return true
	default:
		for i := 3; i < p; i = i + 2 {
			if p%i == 0 {
				return false
			}
		}
		return true
	}
}

/*
//素数を判定するプログラム
func isPrime(p float64) bool {
	count := 0
	//累乗を計算する
	for{
		if count >= 20{
			return true
		}else{
			//pをintに変換
			intp := int(p)

			a := rand.Intn(intp)
			if a==0{
				continue
			}

			//aをfloat64に変換
			flt64a := float64(a)
			num := math.Pow(flt64a,p-1)
			//fmt.Println("aの値:" , a , flt64a)
			//fmt.Println("numの値:",flt64a,"^",p-1 ,"=" , num)

			intnum := int(num)
			modmath := intnum % intp

			//fmt.Println("modmathの値:", intnum , "%" , intp , "=" , modmath)
			if modmath==1{
				count=count+1
			}else{
				return false
			}
		}
	}
}
*/

//ランダムに0~Mまでの素数を作成する関数
func makePrime(M int) int {
	for {
		t := rand.Intn(M)
		if isprime1(t) {
			return t
		}
	}
}

//カギを作成する関数
func makeN(M int) (int, int, int) {
	p := makePrime(M)
	q := makePrime(M)

	N := p * q
	return N, p, q
}

//情報逆元を求めるプログラム
func makegyaku(l, e int) int {
	A := l
	B := e
	a := 0
	b := 1
	for {
		if A%B == 0 {
			return b
		}
		q := A / B
		T := B
		B = A % B
		A = T
		TT := b
		b = a - q*b
		a = TT
	}
}

//最大公約数を作る関数
func gcd(a, b int) int {
	if b == 0 {
		return a
	}
	return gcd(b, a%b)
}

//最小公倍数を作る関数
func lcm(a, b int) int {
	l := a * b / gcd(a, b)
	return l
}

//暗号化
func enc(M, e, N int) int {
	c := math.Pow(float64(M), float64(e))
	c = float64(int(c) % N)
	return int(c)
}

//復号化
func dec(c, d, N int) int {
	M2 := math.Pow(float64(c), float64(d))
	M2 = float64(int(M2) % N)
	return int(M2)
}

func main() {
	rand.Seed(time.Now().Unix())

	//暗号化
	http.HandleFunc("/Enc", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("暗号化を始めます")

		intnum, err := strconv.Atoi(r.FormValue("num"))
		fmt.Println("平文M:",intnum,"を暗号化します")

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		//暗号化
		//N=p*qとなる素数p,qとNを作成
		N, p, q := makeN(intnum)

		//p-1とq-1の最小公倍数lを作成
		l := lcm(p-1, q-1)
		//lとeの最大公約数が1となるeを作成
		e := l - 1
		//情報逆元となるdを計算
		d := makegyaku(l, e)
		fmt.Println("秘密鍵d=", d)
		fmt.Println("公開鍵N=", N)
		fmt.Println("公開鍵e=", e)

		//ここでは公開鍵…N,e 秘密鍵…dとなっている
		angou := enc(intnum, e, N)
		Enctmpl.Execute(w, angou)
		fmt.Println("暗号文=", angou)
		fmt.Println("暗号化を終了します")
	})

	//復号化
	http.HandleFunc("/Dec", func(w http.ResponseWriter, r *http.Request) {

		//http://localhost:8080/Dec?num=1&N=1&e=1&d=1

		fmt.Println("復号化を始めます")
		num := r.FormValue("num")
		N := r.FormValue("N")
		d := r.FormValue("d")
		intN, _ := strconv.Atoi(N)
		intd, _ := strconv.Atoi(d)
		intnum, _ := strconv.Atoi(num)
		fmt.Println(num, N, d)

		hukugou := dec(intnum, intd, intN)
		fmt.Println("hukugou=",hukugou)

		Dectmpl.Execute(w, hukugou)
	})

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
実行例
http://localhost:8080/Enc?num=10
秘密鍵d= -1
公開鍵N= 35
公開鍵e= 11
暗号文= 5

http://localhost:8080/Dec?num=5&N=35&e=11&d=-1

 */
