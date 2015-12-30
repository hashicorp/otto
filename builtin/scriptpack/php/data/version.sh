# php_version_list lists the available PHP versions.
#
# This outputs the result as an array that can be `eval` back into bash.
php_version_list() {
  list=(`apt-cache show php5 | grep Version`)
  for i in "${list[@]}"; do
    if [[ $i == "Version"* ]]; then
      list=(${list[@]/$i})
    fi
  done

  echo "${list[@]}"
}
