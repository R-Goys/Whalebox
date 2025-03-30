package nsenter

/*
#define _GNU_SOURCE
#include "errno.h"
#include "string.h"
#include "stdlib.h"
#include "stdio.h"
#include "sched.h"
#include "fcntl.h"
#include "unistd.h"

__attribute__((constructor)) void enter_namespace(void) {
	char *whalebox_pid;
	whalebox_pid = getenv("whalebox_pid");
	if (whalebox_pid) {
		//fprintf(stdout, "got WHALEBOX_PID: %s\n", whalebox_pid);
	} else {
		//fprintf(stderr, "WHALEBOX_PID not set\n");
		return;
	}
	char *whalebox_cmd;
	whalebox_cmd = getenv("whalebox_cmd");
	if (whalebox_cmd) {
		//fprintf(stdout, "got WHALEBOX_CMD: %s\n", whalebox_cmd);
	} else {
		//fprintf(stdout, "WHALEBOX_CMD not set\n");
		return;
	}
	int i;
	char nspath[1024];
	char *namespace[] = {"net", "pid", "uts", "ipc","mnt"};
	for (i = 0; i < 5; i ++) {
		sprintf(nspath, "/proc/%s/ns/%s", whalebox_pid, namespace[i]);
		int fd = open(nspath, O_RDONLY);
		if (setns(fd, 0) == -1) {
			//fprintf(stderr, "failed to enter %s namespace: %s\n", namespace[i], strerror(errno));
		} else {
			//fprintf(stdout, "entered %s namespace\n", namespace[i]);
		}
		close(fd);
	}
	int res = system(whalebox_cmd);
	exit(0);
	return;
}
*/
import "C"
