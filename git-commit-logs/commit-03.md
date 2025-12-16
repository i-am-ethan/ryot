# commit-03.md

## URL
https://git.kernel.org/pub/scm/git/git.git/commit/?id=e497ea2a9b6c378f01d092c210af20cbee762475

## Content
author	Linus Torvalds <torvalds@ppc970.osdl.org>	2005-04-07 21:03:28 -0700
committer	Linus Torvalds <torvalds@ppc970.osdl.org>	2005-04-07 21:03:28 -0700
commit	e497ea2a9b6c378f01d092c210af20cbee762475 (patch)
tree	4cec02b6467ecc33662b8b691c604d1a80ac4a1b
parent	8bc9a0c769ac1df7820f2dbf8f7b7d64835e3c68 (diff)
download	git-e497ea2a9b6c378f01d092c210af20cbee762475.tar.gz
Make read-tree actually unpack the whole tree.
I needed this to make a "sparse" archive conversion from my old
BitKeeper tree data. The scripts to do the conversion are just
incredibly ugly, but they seem to validate the notion that you
can actually use this silly 'git' thing to save your history in.
Diffstat
-rw-r--r--	read-tree.c	42	
1 files changed, 41 insertions, 1 deletions
diff --git a/read-tree.c b/read-tree.c
index 997343e133..3785e25975 100644
--- a/read-tree.c
+++ b/read-tree.c
@@ -5,6 +5,32 @@
  */
 #include "cache.h"
 
+static void create_directories(const char *path)
+{
+	int len = strlen(path);
+	char *buf = malloc(len + 1);
+	const char *slash = path;
+
+	while ((slash = strchr(slash+1, '/')) != NULL) {
+		len = slash - path;
+		memcpy(buf, path, len);
+		buf[len] = 0;
+		mkdir(buf, 0700);
+	}
+}
+
+static int create_file(const char *path)
+{
+	int fd = open(path, O_WRONLY | O_TRUNC | O_CREAT, 0600);
+	if (fd < 0) {
+		if (errno == ENOENT) {
+			create_directories(path);
+			fd = open(path, O_WRONLY | O_TRUNC | O_CREAT, 0600);
+		}
+	}
+	return fd;
+}
+
 static int unpack(unsigned char *sha1)
 {
 	void *buffer;
@@ -20,12 +46,26 @@ static int unpack(unsigned char *sha1)
 		int len = strlen(buffer)+1;
 		unsigned char *sha1 = buffer + len;
 		char *path = strchr(buffer, ' ')+1;
+		char *data;
+		unsigned long filesize;
 		unsigned int mode;
+		int fd;
+
 		if (size < len + 20 || sscanf(buffer, "%o", &mode) != 1)
 			usage("corrupt 'tree' file");
 		buffer = sha1 + 20;
 		size -= len + 20;
-		printf("%o %s (%s)\n", mode, path, sha1_to_hex(sha1));
+		data = read_sha1_file(sha1, type, &filesize);
+		if (!data || strcmp(type, "blob"))
+			usage("tree file refers to bad file data");
+		fd = create_file(path);
+		if (fd < 0)
+			usage("unable to create file");
+		if (write(fd, data, filesize) != filesize)
+			usage("unable to write file");
+		fchmod(fd, mode);
+		close(fd);
+		free(data);
 	}
 	return 0;
 }