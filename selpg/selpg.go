package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

//selpg的参数
type selpg_args struct {
	start_page  int
	end_page    int
	in_filename string
	print_dest  string
	page_len    int
	page_type   int
}

var progname string //程序名
var sa selpg_args   //当前输入的参数
//var argcount int    //参数个数

const INBUFSIZ = 16 * 1024

func main() {
	flag.IntVar(&sa.start_page, "s", -1, "The start page")
	flag.IntVar(&sa.end_page, "e", -1, "The end page")
	flag.IntVar(&sa.page_len, "l", 72, "The length of the page")
	flag.StringVar(&sa.print_dest, "d", "", "The destination to print")
	f_flag := flag.Bool("f", false, "")
	flag.Parse()

	if *f_flag {
		sa.page_type = 'f'
		sa.page_len = -1
	} else {
		sa.page_type = 'l' // page_type default True
	}
	sa.in_filename = ""
	progname = "selpg"
	//sa.start_page = -1
	//sa.end_page = -1
	//sa.page_len = 72
	//sa.page_type = 'l'
	//sa.print_dest = ""
	//args := os.Args
	//argcount = len(args)
	//process_args(args)
	if flag.NArg() == 1 {
		sa.in_filename = flag.Arg(0)
	}
	validate_args(sa, flag.NArg())
	process_input()
}

func Usage() {
	//fmt.Fprintf(os.Stderr, "\nUSAGE: %s -sstart_page -eend_page [ -f | -llines_per_page ] [ -ddest ] [ in_filename ]\n", progname)
	fmt.Fprintf(os.Stderr, "\nUSAGE: %s --s start_page --e end_page [ --f | --l lines_per_page ] [ --d dest ] [ in_filename ]\n", progname)
}

func validate_args(sa selpg_args, rest int) { // 检验输入参数是否合法，rest为剩余的参数数目
	if rest > 1 {
		fmt.Fprintf(os.Stderr, "./selpg: too much arguments\n")
		Usage()
		os.Exit(1)
	}
	if sa.start_page <= 0 || sa.end_page <= 0 || sa.end_page < sa.start_page {
		fmt.Fprintf(os.Stderr, "./selpg: Invalid start, end page or line number")
		Usage()
		os.Exit(1)
	}
	if sa.page_type == 'f' && sa.page_len != 72 {
		fmt.Fprintf(os.Stderr, "./selpg: Conflict flags: -f and -l")
		Usage()
		os.Exit(1)
	}
}
func process_input() {
	var fin *os.File
	var fout *os.File
	var line_ctr, page_ctr int
	var inpipe io.WriteCloser
	var err error

	if sa.in_filename == "" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(sa.in_filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: ./selpg: could not open input file \"%s\"\n", progname, sa.in_filename)
			Usage()
			os.Exit(1)
		}
		defer fin.Close()
	}
	if sa.print_dest == "" {
		fout = os.Stdout
	} else {
		cmd := exec.Command("lp", "-d", sa.print_dest)
		inpipe, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n", progname, sa.print_dest)
			Usage()
			os.Exit(1)
		}
		defer inpipe.Close()
		cmd.Stdout = fout
		cmd.Start()
	}
	if sa.page_type == 'l' {
		line_ctr = 0
		page_ctr = 1
		line := bufio.NewScanner(fin)
		for line.Scan() {
			if page_ctr >= sa.start_page && page_ctr <= sa.end_page {
				fout.Write([]byte(line.Text() + "\n"))
				if sa.print_dest != "" {
					inpipe.Write([]byte(line.Text() + "\n"))
				}
			}
			line_ctr++
			if line_ctr == sa.page_len {
				page_ctr++
				line_ctr = 0
			}
		}
	} else {
		reader := bufio.NewReader(fin)
		for {
			pageContent, err := reader.ReadString('\f')
			if err != nil || err == io.EOF {
				if err == io.EOF {
					if page_ctr >= sa.start_page && page_ctr <= sa.end_page {
						fmt.Fprintf(fout, "%s", pageContent)
					}
				}
				break
			}
			pageContent = strings.Replace(pageContent, "\f", "", -1)
			if page_ctr >= sa.start_page && page_ctr <= sa.end_page {
				fmt.Fprintf(fout, "%s\n", pageContent)
			}
			page_ctr++
		}
	}
	if page_ctr < sa.start_page {
		fmt.Fprintf(os.Stderr, "./selpg:  start_page (%d) greater than total pages (%d), less output than expected\n", sa.start_page, page_ctr)
	} else if page_ctr < sa.end_page {
		fmt.Fprintf(os.Stderr, "./selpg:  end_page (%d) greater than total pages (%d), less output than expected\n", sa.end_page, page_ctr)
	}

}

//非pflag的参数处理
func process_args(args []string) {
	var argno, i int
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		Usage()
		os.Exit(1)
	}
	/* handle mandatory args first */
	/* handle 1st arg - start page */
	if args[1][0:2] != "-s" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstart_page\n", progname)
		Usage()
		os.Exit(1)
	}
	i, _ = strconv.Atoi(strings.Trim(args[1], "-s"))
	if i < 1 {
		fmt.Fprintf(os.Stderr, "%s: invalid start page %d\n", progname, i)
		Usage()
		os.Exit(1)
	}
	sa.start_page = i
	/* handle 2nd arg - end page */
	if args[2][0:2] != "-e" {
		fmt.Fprintf(os.Stderr, "%s: 2nd arg should be -eend_page\n", progname)
		Usage()
		os.Exit(1)
	}
	i, _ = strconv.Atoi(strings.Trim(args[2], "-e"))
	if i < 1 || i < sa.start_page {
		fmt.Fprintf(os.Stderr, "%s: invalid end page %d\n", progname, i)
		Usage()
		os.Exit(1)
	}
	sa.end_page = i
	/* now handle optional args */
	argno = 3
	for {
		if argno > len(args)-1 || args[argno][0] != '-' {
			break
		}
		switch args[argno][1] {
		case 'l':
			//获取一页的长度
			i, _ := strconv.Atoi(args[argno][2:])
			if i < 1 {
				fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, i)
				Usage()
				os.Exit(1)
			}
			sa.page_len = i
			argno++
		case 'f':
			if args[argno] != "-f" {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				Usage()
				os.Exit(1)
			}
			sa.page_type = 'f'
			argno++
		case 'd':
			if args[argno] == "-d" {
				fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination\n", progname)
				Usage()
				os.Exit(1)
			}
			sa.print_dest = args[argno][2:]
			argno++
		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option", progname)
			Usage()
			os.Exit(1)
		}
	}
	if argno <= len(args)-1 {
		sa.in_filename = args[argno]
		_, err := os.Stat(sa.in_filename)
		if err != nil {
			/* check if file exists */
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n",
					progname, sa.in_filename)
				os.Exit(1)
			}
			/* check if file is readable */
			if os.IsPermission(err) {
				fmt.Fprintf(os.Stderr, "%s: input file \"%s\" exists but cannot be read\n",
					progname, sa.in_filename)
				os.Exit(1)
			}
		}
	}
}
