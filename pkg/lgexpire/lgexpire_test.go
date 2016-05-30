package lgexpire

import (
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"
)

func TestLogFileNameParsing(t *testing.T) {
	names := []string{
		"very.cool.program.raspberrypi.unknownuser.log.INFO.20160521-235713.736",
		"very-cool-program.coolhost.cooluser.log.ERROR.20160522-103338.8664",
		"very-cool-program.coolhost.cooluser.log.FATAL.20160522-103338.8664",
		"very-cool-program.coolhost.cooluser.log.INFO.20160522-103338.8664",
		"very-cool-program.coolhost.cooluser.log.INFO.20160522-103346.8732",
		"very-cool-program.coolhost.cooluser.log.INFO.20160522-103402.9194",
		"very-cool-program.coolhost.cooluser.log.INFO.20160522-103406.9305",
		"very-cool-program.coolhost.cooluser.log.INFO.20160522-103414.9416",
		"very-cool-program.coolhost.cooluser.log.WARNING.20160522-103338.8664.gz",
		"very-cool-program.coolhost.cooluser.log.WARNING.20160522-103338.8664.gz",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160521-235742.757",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160522-103555.757",
		"dino-catcher.raspberrypi.unknownuser.log.WARNING.20160521-235742.757",
		"dino-catcher.raspberrypi.unknownuser.log.WARNING.20160522-103555.757",
		"disco-central.disco.root.log.ERROR.20160302-004210.28171",
		"disco-central.disco.root.log.ERROR.20160302-004512.28873",
		"disco-central.disco.root.log.ERROR.20160302-004519.28883",
		"disco-central.disco.root.log.ERROR.20160302-004522.28898",
		"disco-central.disco.root.log.FATAL.20160302-004210.28171",
		"disco-central.disco.root.log.FATAL.20160302-004512.28873",
		"disco-central.disco.root.log.FATAL.20160302-004519.28883",
		"disco-central.disco.root.log.FATAL.20160302-004522.28898",
		"disco-central.disco.root.log.INFO.20160302-004210.28171",
		"disco-central.disco.root.log.INFO.20160302-004512.28873",
		"disco-central.disco.root.log.INFO.20160302-004519.28883",
		"disco-central.disco.root.log.INFO.20160302-004522.28898",
		"disco-central.disco.root.log.WARNING.20160302-004210.28171",
		"disco-central.disco.root.log.WARNING.20160302-004512.28873",
		"disco-central.disco.root.log.WARNING.20160302-004519.28883",
		"disco-central.disco.root.log.WARNING.20160302-004522.28898",
		"disco-dance-server.disco.disco-dance-server.log.INFO.20160217-202201.4385",
		"disco-dance-server.disco.disco-dance-server.log.INFO.20160302-002715.24192",
		"disco-dance-server.disco.disco-dance-server.log.INFO.20160510-211059.18831",
	}
	for _, name := range names {
		lf, err := parseLogFileName(name)
		if err != nil {
			t.Fatalf("test of %s failed: %v", name, err)
		}
		_ = lf
	}

	invalidNames := []string{
		"very.cool.program.raspberrypi.unknownuser.log.INFA.20160521-235713.736",
		"very-cool-program.coolhost.cooluser.log.ERROR.20160522-103338",
		"very-cool-program.FATAL",
	}
	for _, name := range invalidNames {
		_, err := parseLogFileName(name)
		if err == nil {
			t.Fatalf("name %s should have failed", name)
		}

	}
}

func TestRotate(t *testing.T) {

	// only one file of each level should remainc
	r := &Expire{
		Programs: []string{"any-central"},
		Rules: []Rule{
			{Count: 1},
			{Count: 2},
		},
	}
	l := newLocalTest(t, r, 1024*1024, infiles)

	l.Run(func(r *Expire) {
		err := r.Run()
		if err != nil {
			t.Fatal(err)
		}
		l.AssertFiles(
			"any-admin_linux_amd64.INFO",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011043.27029",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011048.27039",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011348.27614",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-012702.30128",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-013907.798",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-014322.1669",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-023706.12116",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-025324.15722",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-032215.21939",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-033637.24399",
			"any-admin_linux_amd64.any.root.log.INFO.20160212-025702.12553",
			"any-admin_linux_amd64.any.root.log.INFO.20160216-220156.19842",
			"any-admin_linux_amd64.any.root.log.INFO.20160225-122849.12392",
			"any-admin_linux_amd64.any.root.log.INFO.20160302-172837.31078",
			"any-central.ERROR",
			"any-central.FATAL",
			"any-central.INFO",
			"any-central.WARNING",
			"any-central.any.any-central.log.ERROR.20160403-170319.28594",
			"any-central.any.any-central.log.FATAL.20160217-202203.4384",
			"any-central.any.any-central.log.INFO.20160403-162539.28594",
			"any-central.any.any-central.log.WARNING.20160403-165721.28594",
		)
	})

}
func TestRotate2(t *testing.T) {

	// all files should be older than one second
	r := &Expire{
		Programs: []string{"any-central"},
		Rules: []Rule{
			{Count: 5, Age: time.Second},
		},
	}
	l := newLocalTest(t, r, 1024*1024, infiles)

	l.Run(func(r *Expire) {
		err := r.Run()
		if err != nil {
			t.Fatal(err)
		}
		l.AssertFiles(
			"any-admin_linux_amd64.INFO",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011048.27039",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011348.27614",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-012702.30128",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-013907.798",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-014322.1669",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-023706.12116",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-025324.15722",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-032215.21939",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-033637.24399",
			"any-admin_linux_amd64.any.root.log.INFO.20160212-025702.12553",
			"any-admin_linux_amd64.any.root.log.INFO.20160216-220156.19842",
			"any-admin_linux_amd64.any.root.log.INFO.20160225-122849.12392",
			"any-admin_linux_amd64.any.root.log.INFO.20160302-172837.31078",
			"any-central.ERROR",
			"any-central.FATAL",
			"any-central.INFO",
			"any-central.WARNING",
			"any-central.any.any-central.log.ERROR.20160403-170319.28594",
			"any-central.any.any-central.log.FATAL.20160217-202203.4384",
			"any-central.any.any-central.log.INFO.20160403-162539.28594",
			"any-central.any.any-central.log.WARNING.20160403-165721.28594",
			"any-admin_linux_amd64.any.root.log.INFO.20160211-011043.27029",
		)
	})
}
func TestRotate3(t *testing.T) {

	// all files should be kept
	r := &Expire{
		Programs: []string{"any-central"},
		Rules: []Rule{
			{Count: 500, Age: 24 * 365 * 200 * time.Hour},
		},
	}
	l := newLocalTest(t, r, 1024*1024, infiles)

	l.Run(func(r *Expire) {
		err := r.Run()
		if err != nil {
			t.Fatal(err)
		}
		l.AssertFiles(infiles...)

	})
}
func TestRotate4(t *testing.T) {

	// at least two of each program log should be kept
	r := &Expire{
		Programs: []string{"any-central", "any-admin_linux_amd64"},
		Rules: []Rule{
			{Count: 2, Age: 24 * 365 * 200 * time.Hour},
			{Count: 5},
		},
	}
	l := newLocalTest(t, r, 1024*1024, infiles)

	l.Run(func(r *Expire) {
		err := r.Run()
		if err != nil {
			t.Fatal(err)
		}
		l.AssertFiles(
			"any-admin_linux_amd64.INFO",
			"any-admin_linux_amd64.any.root.log.INFO.20160225-122849.12392",
			"any-admin_linux_amd64.any.root.log.INFO.20160302-172837.31078",
			"any-central.ERROR",
			"any-central.FATAL",
			"any-central.INFO",
			"any-central.WARNING",
			"any-central.any.any-central.log.ERROR.20160319-104215.1587",
			"any-central.any.any-central.log.ERROR.20160403-170319.28594",
			"any-central.any.any-central.log.FATAL.20160212-090515.26564",
			"any-central.any.any-central.log.FATAL.20160217-202203.4384",
			"any-central.any.any-central.log.INFO.20160317-012126.1587",
			"any-central.any.any-central.log.INFO.20160403-162539.28594",
			"any-central.any.any-central.log.WARNING.20160317-030816.1587",
			"any-central.any.any-central.log.WARNING.20160403-165721.28594",
		)
	})
}

func TestRotate5(t *testing.T) {

	infiles := []string{
		"dino-catcher.INFO",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-075157.734",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-075741.767",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-103220.767",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-104136.767",
		"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-121407.1068", // this one was deleted
	}
	r := &Expire{
		Programs: []string{"dino-catcher"},
		Rules:    []Rule{{Count: 4}},
	}
	l := newLocalTest(t, r, 1024*1024, infiles)

	l.Run(func(r *Expire) {
		err := r.Run()
		if err != nil {
			t.Fatal(err)
		}
		l.AssertFiles(
			"dino-catcher.INFO",
			"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-075741.767",
			"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-103220.767",
			"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-104136.767",
			"dino-catcher.raspberrypi.unknownuser.log.INFO.20160525-121407.1068",
		)
	})

}

// CreateFile .
type CreateFile struct {
	name string
	size int
}

func newLocalTest(t *testing.T, r *Expire, size int, files []string) *LocalTest {
	lt := &LocalTest{
		t:           t,
		rotate:      r,
		maxFilesize: size,
		files:       files,
	}
	return lt
}

// LocalTest .
type LocalTest struct {
	t           *testing.T
	rotate      *Expire
	maxFilesize int      // max file size to generate
	files       []string // file names
	tmpdir      string
}

func (l *LocalTest) Run(f func(r *Expire)) {
	tmpdir, err := ioutil.TempDir("", "lgrotate-test")
	if err != nil {
		panic(err)
	}
	l.tmpdir = tmpdir
	// l.T.Parallel()
	defer func() {
		os.RemoveAll(tmpdir)
	}()

	for _, fn := range l.files {
		fullname := filepath.Join(tmpdir, fn)
		err := ioutil.WriteFile(fullname, []byte("content"), 0777)
		if err != nil {
			l.t.Fatal(err)
		}
	}
	l.rotate.LogDir = tmpdir

	f(l.rotate)
}

// AssertFiles fails the test if the local maven resulting repo doenst
// contain exactly the files specified by the path arguments.
func (l *LocalTest) AssertFiles(path ...string) {
	ok, files := l.expectFiles(path...)
	if !ok {
		l.t.Fatalf(
			"unexpected file situation:\n\nfound:\n\n%s\n\nexpected:\n\n%s\n\n",
			strings.Join(files, "\n"),
			strings.Join(path, "\n"),
		)
	}
}

// AssertFiles fails the test if the local maven resulting repo doenst
// contain exactly the files specified by the path arguments.
func (l *LocalTest) AssertNoFiles() {
	ok, files := l.expectFiles("")
	if ok {
		l.t.Fatalf(
			"expeced no files, got: \n %s\n\n",
			strings.Join(files, "\n"),
		)
	}
}

func (l *LocalTest) expectFiles(path ...string) (bool, []string) {

	var files []string
	err := filepath.Walk(l.tmpdir, func(path string, f os.FileInfo, err error) error {
		p, err := filepath.Rel(l.tmpdir, path)
		if err != nil {
			panic(err)
		}
		if !f.IsDir() {
			files = append(files, p)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	sort.Strings(files)
	sort.Strings(path)
	ok := true
	if !reflect.DeepEqual(files, path) {
		ok = false
	}
	return ok, files
}

// shuffle shuffles a slice using the Fisher-Yates algoritm.
func shuffle(slc []string) {
	N := len(slc)
	for i := 0; i < N; i++ {
		// choose index uniformly in [i, N-1]
		r := i + rand.Intn(N-i)
		slc[r], slc[i] = slc[i], slc[r]
	}
}

var infiles = []string{
	"any-admin_linux_amd64.INFO",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-011043.27029",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-011048.27039",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-011348.27614",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-012702.30128",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-013907.798",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-014322.1669",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-023706.12116",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-025324.15722",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-032215.21939",
	"any-admin_linux_amd64.any.root.log.INFO.20160211-033637.24399",
	"any-admin_linux_amd64.any.root.log.INFO.20160212-025702.12553",
	"any-admin_linux_amd64.any.root.log.INFO.20160216-220156.19842",
	"any-admin_linux_amd64.any.root.log.INFO.20160225-122849.12392",
	"any-admin_linux_amd64.any.root.log.INFO.20160302-172837.31078",
	"any-central.ERROR",
	"any-central.FATAL",
	"any-central.INFO",
	"any-central.WARNING",
	"any-central.any.any-central.log.ERROR.20160209-091120.7851",
	"any-central.any.any-central.log.ERROR.20160212-090515.26564",
	"any-central.any.any-central.log.ERROR.20160217-202203.4384",
	"any-central.any.any-central.log.ERROR.20160225-104310.21294",
	"any-central.any.any-central.log.ERROR.20160225-183945.3328",
	"any-central.any.any-central.log.ERROR.20160302-104420.28562",
	"any-central.any.any-central.log.ERROR.20160319-104215.1587",
	"any-central.any.any-central.log.ERROR.20160403-170319.28594",
	"any-central.any.any-central.log.FATAL.20160209-091120.7851",
	"any-central.any.any-central.log.FATAL.20160212-090515.26564",
	"any-central.any.any-central.log.FATAL.20160217-202203.4384",
	"any-central.any.any-central.log.INFO.20160131-022239.17027",
	"any-central.any.any-central.log.INFO.20160131-022442.17717",
	"any-central.any.any-central.log.INFO.20160131-022529.17946",
	"any-central.any.any-central.log.INFO.20160131-023735.21245",
	"any-central.any.any-central.log.INFO.20160131-024015.23179",
	"any-central.any.any-central.log.INFO.20160131-024258.23894",
	"any-central.any.any-central.log.INFO.20160131-024502.24370",
	"any-central.any.any-central.log.INFO.20160131-024554.24870",
	"any-central.any.any-central.log.INFO.20160131-030059.29900",
	"any-central.any.any-central.log.INFO.20160131-030119.30002",
	"any-central.any.any-central.log.INFO.20160131-030329.30746",
	"any-central.any.any-central.log.INFO.20160131-142709.6408",
	"any-central.any.any-central.log.INFO.20160131-161242.6649",
	"any-central.any.any-central.log.INFO.20160204-022533.5060",
	"any-central.any.any-central.log.INFO.20160208-165752.4279",
	"any-central.any.any-central.log.INFO.20160208-182930.7851",
	"any-central.any.any-central.log.INFO.20160209-091121.25120",
	"any-central.any.any-central.log.INFO.20160211-034329.26564",
	"any-central.any.any-central.log.INFO.20160212-090516.25043",
	"any-central.any.any-central.log.INFO.20160216-152057.1041",
	"any-central.any.any-central.log.INFO.20160216-232237.4931",
	"any-central.any.any-central.log.INFO.20160217-202201.4384",
	"any-central.any.any-central.log.INFO.20160217-202204.4731",
	"any-central.any.any-central.log.INFO.20160225-104031.21294",
	"any-central.any.any-central.log.INFO.20160225-111135.27860",
	"any-central.any.any-central.log.INFO.20160225-114745.3328",
	"any-central.any.any-central.log.INFO.20160302-002714.24185",
	"any-central.any.any-central.log.INFO.20160302-004133.28015",
	"any-central.any.any-central.log.INFO.20160302-004335.28562",
	"any-central.any.any-central.log.INFO.20160317-012126.1587",
	"any-central.any.any-central.log.INFO.20160403-162539.28594",
	"any-central.any.any-central.log.WARNING.20160131-082329.30746",
	"any-central.any.any-central.log.WARNING.20160131-143719.6408",
	"any-central.any.any-central.log.WARNING.20160131-162252.6649",
	"any-central.any.any-central.log.WARNING.20160204-024750.5060",
	"any-central.any.any-central.log.WARNING.20160208-210925.7851",
	"any-central.any.any-central.log.WARNING.20160209-201725.25120",
	"any-central.any.any-central.log.WARNING.20160211-075329.26564",
	"any-central.any.any-central.log.WARNING.20160212-140914.25043",
	"any-central.any.any-central.log.WARNING.20160216-174638.1041",
	"any-central.any.any-central.log.WARNING.20160217-004237.4931",
	"any-central.any.any-central.log.WARNING.20160217-202203.4384",
	"any-central.any.any-central.log.WARNING.20160217-202537.4731",
	"any-central.any.any-central.log.WARNING.20160225-104310.21294",
	"any-central.any.any-central.log.WARNING.20160225-122749.3328",
	"any-central.any.any-central.log.WARNING.20160302-004925.28562",
	"any-central.any.any-central.log.WARNING.20160317-030816.1587",
	"any-central.any.any-central.log.WARNING.20160403-165721.28594",
}
