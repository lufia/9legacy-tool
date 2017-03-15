# manage 9legacy patches

## INSTALL

	chmod +x apply cache init install installall list stable uninstall update

## SYNOPSIS

	9legacy/init
	9legacy/update
	9legacy/install {patch-name}
	9legacy/uninstall {patch-name}
	9legacy/list
	9legacy/installall [file...]
	9legacy/stable

	9legacy/apply

	9legacy/cache {patch-name}

## DESCRIPTION

9legacy-tool manages to apply 9legacy's patches to Plan 9 box.
It scrapes 9legacy.org, downloads patches, and updates sources indirectly.

## FILES

* $home/lib/9legacy - 9legacy-tool's working directory

## EXAMPLE

These instructions installs all stable patches into $home/lib/9legacy.

	# initialize working directory
	% 9legacy/init

	# update available patch list from 9legacy.org
	% 9legacy/update

	# install stable patches
	% 9legacy/installall <{9legacy/stable}

And `9legacy/apply` apply installed patches onto the system by constructing namespace.

	% 9legacy/apply

## BUGS

After `9legacy/apply`, writing a new file to subdirectory just under the root will redirect to working directory.
For example, writing a file to /amd64/init redirects to $home/lib/9legacy/plan9/amd64/init.
