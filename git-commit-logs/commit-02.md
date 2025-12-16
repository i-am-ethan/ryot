# commit-02.md

## URL
https://git.kernel.org/pub/scm/git/git.git/commit/?id=8bc9a0c769ac1df7820f2dbf8f7b7d64835e3c68

## Content
author	Linus Torvalds <torvalds@ppc970.osdl.org>	2005-04-07 15:16:10 -0700
committer	Linus Torvalds <torvalds@ppc970.osdl.org>	2005-04-07 15:16:10 -0700
commit	8bc9a0c769ac1df7820f2dbf8f7b7d64835e3c68 (patch)
tree	4f856748145b27bc669a98cd4d66574dfd01b083
parent	e83c5163316f89bfbde7d9ab23ca2e25604af290 (diff)
download	git-8bc9a0c769ac1df7820f2dbf8f7b7d64835e3c68.tar.gz
Add copyright notices.
The tool interface sucks (especially "committing" information, which is just
me doing everything by hand from the command line), but I think this is in
theory actually a viable way of describing the world. So copyright it.
Diffstat
-rw-r--r--	cat-file.c	5	
-rw-r--r--	commit-tree.c	5	
-rw-r--r--	init-db.c	5	
-rw-r--r--	read-cache.c	5	
-rw-r--r--	read-tree.c	5	
-rw-r--r--	show-diff.c	5	
-rw-r--r--	update-cache.c	5	
-rw-r--r--	write-tree.c	5	
8 files changed, 40 insertions, 0 deletions
diff --git a/cat-file.c b/cat-file.c
index 74a0a234dd..d8f0121e49 100644
--- a/cat-file.c
+++ b/cat-file.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 int main(int argc, char **argv)
diff --git a/commit-tree.c b/commit-tree.c
index 840307af0c..0d7852f7f2 100644
--- a/commit-tree.c
+++ b/commit-tree.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 #include <pwd.h>
diff --git a/init-db.c b/init-db.c
index 25dc13fe10..bc65086368 100644
--- a/init-db.c
+++ b/init-db.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 int main(int argc, char **argv)
diff --git a/read-cache.c b/read-cache.c
index c924a6e0fc..b151981c1b 100644
--- a/read-cache.c
+++ b/read-cache.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 const char *sha1_file_directory = NULL;
diff --git a/read-tree.c b/read-tree.c
index 1b47742d8c..997343e133 100644
--- a/read-tree.c
+++ b/read-tree.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 static int unpack(unsigned char *sha1)
diff --git a/show-diff.c b/show-diff.c
index b8522886a1..e7aa60b4ab 100644
--- a/show-diff.c
+++ b/show-diff.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 #define MTIME_CHANGED	0x0001
diff --git a/update-cache.c b/update-cache.c
index 5085a5cb53..0bcd5d4081 100644
--- a/update-cache.c
+++ b/update-cache.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 static int cache_name_compare(const char *name1, int len1, const char *name2, int len2)
diff --git a/write-tree.c b/write-tree.c
index 921f981353..012eeeff1b 100644
--- a/write-tree.c
+++ b/write-tree.c
@@ -1,3 +1,8 @@
+/*
+ * GIT - The information manager from hell
+ *
+ * Copyright (C) Linus Torvalds, 2005
+ */
 #include "cache.h"
 
 static int check_valid_sha1(unsigned char *sha1)
