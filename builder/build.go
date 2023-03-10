package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	myproto "taskdev/proto"
	"taskdev/util"
	"time"
)

func (b *builder) addBuild(key string) {
	b.buildQueueLock.Lock()
	defer b.buildQueueLock.Unlock()
	b.buildQueue[key] = key
}
func (b *builder) delBuild(key string) {
	b.buildQueueLock.Lock()
	defer b.buildQueueLock.Unlock()
	if _, ok := b.buildQueue[key]; ok {
		delete(b.buildQueue, key)
	}
}
func (b *builder) checkBuild(key string) bool {
	b.buildQueueLock.Lock()
	defer b.buildQueueLock.Unlock()
	_, ok := b.buildQueue[key]
	return ok
}

func (b *builder) getSvnCode(m *myproto.TaskBuildReq_TaskBuild, buildDir string, buildLog string) (string, string, error) {
	version := m.GetVersion()
	repos := m.GetRepos()
	i := strings.LastIndex(repos, "/")
	if i != -1 && i != len(repos)-1 {
		buildDir = buildDir + repos[i+1:]
	}
	b.log.Debug("build DIR = %s\n", b.buildDir)
	b.log.Debug("task build dir: %s\n", buildDir)

	var svnArgs []string
	svnDir := buildDir + "/.svn"
	if _, err := os.Stat(svnDir); err == nil {
		svnArgs = append(svnArgs, "update")
		svnArgs = append(svnArgs, buildDir)

	} else {
		user := m.GetUser()
		pass := m.GetPass()

		svnArgs = append(svnArgs, "checkout")
		svnArgs = append(svnArgs, "--username")
		svnArgs = append(svnArgs, user)
		svnArgs = append(svnArgs, "--password")
		svnArgs = append(svnArgs, pass)
		svnArgs = append(svnArgs, fmt.Sprintf("%s", repos))
		svnArgs = append(svnArgs, buildDir)
	}

	b.log.Debug("execute command %s\n", fmt.Sprintf("svn %s", svnArgs))
	buildLog = fmt.Sprintf("%sBUILD: svn %s\n", buildLog, svnArgs)

	cmd := exec.Command("svn", svnArgs...)
	out, err := cmd.CombinedOutput()
	b.log.Info("execute out:\n%s\n", out)
	buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, out)
	if err != nil {
		b.log.Error("execute error, err = %s\n", err)
		buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, err)
		return "", buildLog, err

	}
	prjDir := m.GetReposDir()
	if len(prjDir) > 0 {
		if !strings.HasSuffix(buildDir, string(filepath.Separator)) &&
			!strings.HasPrefix(prjDir, string(filepath.Separator)) {
			prjDir = buildDir + string(filepath.Separator) + prjDir
		} else {
			prjDir = buildDir + prjDir
		}
		svnArgs = nil
		svnArgs = append(svnArgs, "update")
		if len(version) == 0 {
			svnArgs = append(svnArgs, prjDir)
		} else {
			svnArgs = append(svnArgs, fmt.Sprintf("-r%s", version))
			svnArgs = append(svnArgs, prjDir)
		}
		buildLog = fmt.Sprintf("%sBUILD: svn %s\n", buildLog, svnArgs)
		cmd := exec.Command("svn", svnArgs...)
		out, err := cmd.CombinedOutput()
		b.log.Info("execute out:\n%s\n", out)
		buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, out)
		if err != nil {
			b.log.Error("execute error, err = %s\n", err)
			buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, err)
			return "", buildLog, err

		}
	}
	return buildLog, buildDir, nil
}

func (b *builder) build(m *myproto.TaskBuildReq_TaskBuild) {

	name := m.GetName()
	groupName := m.GetGroupName()
	version := m.GetVersion()

	b.addBuild(fmt.Sprintf("%s-%s", groupName, name))

	outbin := m.GetOutBin()

	buildLog := ""
	buildLogDir := ""
	buildScriptDir := ""

	st := &buildState{
		groupname: groupName,
		name:      name,
		state:     0,
	}

	defer func() {
		buildLog = fmt.Sprintf("%sBUILD: DONE!\n", buildLog)
		n := time.Now()
		fileName := fmt.Sprintf("build-%04d%02d%02d%02d%02d%02d.log",
			n.Year(), n.Month(), n.Day(), n.Hour(), n.Minute(), n.Second())
		prefix := groupName
		if len(name) > 0 {
			if len(groupName) > 0 {
				prefix = prefix + "-" + name
			} else {
				prefix = name
			}
			fileName = prefix + "-" + fileName
		} else {
			if len(groupName) > 0 {
				fileName = prefix + "-" + fileName
			}
		}
		fileName = buildLogDir + fileName
		ioutil.WriteFile(fileName, ([]byte)(buildLog), 0644)

		st.log = fileName

		b.db.insert(st)

		b.delBuild(fmt.Sprintf("%s-%s", groupName, name))

	}()

	b.log.Info("start to build task, g=%s, n=%s, v=%s\n",
		groupName, name, version)

	buildDir := b.buildDir + groupName
	if len(groupName) != 0 {
		if !strings.HasSuffix(buildDir, string(filepath.Separator)) {
			buildDir = buildDir + string(filepath.Separator)
		}
		if _, err := os.Stat(buildDir); err != nil {
			os.Mkdir(buildDir, 0777)
		}
	}
	buildDir = buildDir + name
	if len(name) != 0 {
		if !strings.HasSuffix(buildDir, string(filepath.Separator)) {
			buildDir = buildDir + string(filepath.Separator)
		}
		if _, err := os.Stat(buildDir); err != nil {
			os.Mkdir(buildDir, 0777)
		}
	}
	buildLogDir = buildDir + "buildlog" + string(filepath.Separator)
	if _, err := os.Stat(buildLogDir); err != nil {
		os.Mkdir(buildLogDir, 0777)
	}

	buildScriptDir = buildDir + "buildscript" + string(filepath.Separator)
	if _, err := os.Stat(buildScriptDir); err != nil {
		os.Mkdir(buildScriptDir, 0777)
	}

	for _, v := range outbin {
		if len(st.outbin) != 0 {
			st.outbin = st.outbin + ";"
		}
		st.outbin = buildDir + st.outbin + v
	}

	vcs := m.GetVcs()
	switch {
	case vcs == util.KVCSSvn:
		{
			var err error
			buildLog, buildDir, err = b.getSvnCode(m, buildDir, buildLog)
			if err != nil {
				b.log.Debug("getSvnCode failed, err = %s\n", err)
				buildLog = fmt.Sprintf("%sBUILD: FAILED %s\n", buildLog, err)
				st.state = buildStateBuildFailed
				return
			}
			b.log.Debug("get svn code done\n")
		}
	}
	preBuildScript := m.GetPreBuildScript()
	buildScript := m.GetBuildScript()
	postBuildScript := m.GetPostBuildScript()

	buildScriptCommonName := ""
	if len(groupName) > 0 {
		buildScriptCommonName = groupName
	}
	if len(name) > 0 {
		if len(buildScriptCommonName) > 0 {
			buildScriptCommonName = buildScriptCommonName + "-" + name
		} else {
			buildScriptCommonName = name
		}
	}
	buildCmmNameLen := len(buildScriptCommonName)
	preBuildName := buildScriptCommonName

	if len(preBuildScript) > 0 {
		b.log.Debug("start prebuild\n")
		if buildCmmNameLen > 0 {
			preBuildName = buildScriptDir + preBuildName + "-prebuild.sh"
		} else {
			preBuildName = buildScriptDir + "prebuild.sh"
		}
		err := ioutil.WriteFile(preBuildName, []byte(preBuildScript), 0744)
		if err != nil {
			b.log.Error("create prebuild script failed.\n")
			buildLog = fmt.Sprintf("%sBUILD: FAILED %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		buildLog = fmt.Sprintf("%sBUILD: prebuild %s\n", buildLog, preBuildName)
		cmd := exec.Command(preBuildName)
		cmd.Dir = buildDir
		out, err := cmd.CombinedOutput()
		b.log.Info("prebuild out:\n%s\n", out)
		buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, out)
		if err != nil {
			b.log.Error("build error, err = %s\n", err)
			buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		b.log.Debug("prebuild done\n")

	}
	buildName := buildScriptCommonName
	if len(buildScript) > 0 {
		b.log.Debug("start build\n")

		if buildCmmNameLen > 0 {
			buildName = buildScriptDir + buildName + "-build.sh"
		} else {
			buildName = buildScriptDir + "build.sh"
		}
		err := ioutil.WriteFile(buildName, []byte(buildScript), 0744)
		if err != nil {
			b.log.Error("create prebuild script failed.\n")
			buildLog = fmt.Sprintf("%sBUILD: FAILED %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		buildLog = fmt.Sprintf("%sBUILD: build %s\n", buildLog, buildName)
		cmd := exec.Command(buildName)
		cmd.Dir = buildDir
		out, err := cmd.CombinedOutput()
		b.log.Info("build out:\n%s\n", out)
		buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, out)
		if err != nil {
			b.log.Error("build error, err = %s\n", err)
			buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		b.log.Debug("build done\n")

	}
	postBuildName := buildScriptCommonName
	if len(postBuildScript) > 0 {
		b.log.Debug("start postbuild\n")

		if buildCmmNameLen > 0 {
			postBuildName = buildScriptDir + postBuildName + "-postbuild.sh"
		} else {
			postBuildName = buildScriptDir + "postbuild.sh"
		}
		err := ioutil.WriteFile(postBuildName, []byte(postBuildScript), 0744)
		if err != nil {
			b.log.Error("create prebuild script failed.\n")
			buildLog = fmt.Sprintf("%sBUILD: FAILED %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		buildLog = fmt.Sprintf("%sBUILD: postbuild %s\n", buildLog, postBuildName)
		cmd := exec.Command(postBuildName)
		cmd.Dir = buildDir
		out, err := cmd.CombinedOutput()
		b.log.Info("build out:\n%s\n", out)
		buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, out)
		if err != nil {
			b.log.Error("build error, err = %s\n", err)
			buildLog = fmt.Sprintf("%sBUILD: %s\n", buildLog, err)
			st.state = buildStateBuildFailed

			return
		}
		b.log.Debug("postbuild done\n")

	}

	st.state = buildStateBuildSuccess
	b.log.Debug("all build done\n")

}
