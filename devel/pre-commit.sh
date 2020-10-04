#!/usr/bin/sh

yield()
{
	echo "$*" >&2
}

error_msg()
{
	yield "ERROR: $*"
}

die()
{
	err=$?
	[ $err -eq 0 ] && err=127
	error_msg "$@"
	return $err
}

if git rev-parse --verify HEAD >/dev/null 2>&1
then
	against="HEAD"
else
	# Initial commit: diff against an empty tree object
	against="$( git hash-object -t tree /dev/null )"
fi

# If you want to allow non-ASCII filenames set this variable to true.
allownonascii="$( git config --type=bool hooks.allownonascii )"

# Redirect output to stderr.
exec 1>&2

# Cross platform projects tend to avoid non-ASCII filenames; prevent
# them from being added to the repository. We exploit the fact that the
# printable range starts at the space character and ends with tilde.
if [ "${allownonascii}" != "true" ] &&
	# Note that the use of brackets around a tr range is ok here, (it's
	# even required, for portability to Solaris 10's /usr/bin/tr), since
	# the square bracket bytes happen to fall in the designated range.
	test "$( git diff --cached --name-only --diff-filter=A -z "${against}" \
	         | LC_ALL=C tr -d '[ -~]\0' \
		     | wc -c)" != 0
then
	cat <<\EOF
Error: Attempt to add a non-ASCII file name.
This can cause problems if you want to work with people on other platforms.
To be portable it is advisable to rename the file.
If you know what you are doing you can disable this check using:
  git config hooks.allownonascii true
EOF
	exit 1
fi

# If there are whitespace errors, print the offending file names and fail.
git diff-index --check --cached "${against}" -- || die "Remove the whitespaces and commit again"

files=""
for item in $( git diff-index --cached --name-only "${against}" 2>/dev/null )
do
	files="${files} ${item}"
done

if [ "${files}" != "" ]; then
    echo "Linting:${files}"
    # shellcheck disable=SC2086
    ./devel/lint.sh ${files}
    exit $?
else
    exit 0
fi

