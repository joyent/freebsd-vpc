# bash completion for vpc                                  -*- shell-script -*-

__debug()
{
    if [[ -n ${BASH_COMP_DEBUG_FILE} ]]; then
        echo "$*" >> "${BASH_COMP_DEBUG_FILE}"
    fi
}

# Homebrew on Macs have version 1.3 of bash-completion which doesn't include
# _init_completion. This is a very minimal version of that function.
__my_init_completion()
{
    COMPREPLY=()
    _get_comp_words_by_ref "$@" cur prev words cword
}

__index_of_word()
{
    local w word=$1
    shift
    index=0
    for w in "$@"; do
        [[ $w = "$word" ]] && return
        index=$((index+1))
    done
    index=-1
}

__contains_word()
{
    local w word=$1; shift
    for w in "$@"; do
        [[ $w = "$word" ]] && return
    done
    return 1
}

__handle_reply()
{
    __debug "${FUNCNAME[0]}"
    case $cur in
        -*)
            if [[ $(type -t compopt) = "builtin" ]]; then
                compopt -o nospace
            fi
            local allflags
            if [ ${#must_have_one_flag[@]} -ne 0 ]; then
                allflags=("${must_have_one_flag[@]}")
            else
                allflags=("${flags[*]} ${two_word_flags[*]}")
            fi
            COMPREPLY=( $(compgen -W "${allflags[*]}" -- "$cur") )
            if [[ $(type -t compopt) = "builtin" ]]; then
                [[ "${COMPREPLY[0]}" == *= ]] || compopt +o nospace
            fi

            # complete after --flag=abc
            if [[ $cur == *=* ]]; then
                if [[ $(type -t compopt) = "builtin" ]]; then
                    compopt +o nospace
                fi

                local index flag
                flag="${cur%%=*}"
                __index_of_word "${flag}" "${flags_with_completion[@]}"
                COMPREPLY=()
                if [[ ${index} -ge 0 ]]; then
                    PREFIX=""
                    cur="${cur#*=}"
                    ${flags_completion[${index}]}
                    if [ -n "${ZSH_VERSION}" ]; then
                        # zsh completion needs --flag= prefix
                        eval "COMPREPLY=( \"\${COMPREPLY[@]/#/${flag}=}\" )"
                    fi
                fi
            fi
            return 0;
            ;;
    esac

    # check if we are handling a flag with special work handling
    local index
    __index_of_word "${prev}" "${flags_with_completion[@]}"
    if [[ ${index} -ge 0 ]]; then
        ${flags_completion[${index}]}
        return
    fi

    # we are parsing a flag and don't have a special handler, no completion
    if [[ ${cur} != "${words[cword]}" ]]; then
        return
    fi

    local completions
    completions=("${commands[@]}")
    if [[ ${#must_have_one_noun[@]} -ne 0 ]]; then
        completions=("${must_have_one_noun[@]}")
    fi
    if [[ ${#must_have_one_flag[@]} -ne 0 ]]; then
        completions+=("${must_have_one_flag[@]}")
    fi
    COMPREPLY=( $(compgen -W "${completions[*]}" -- "$cur") )

    if [[ ${#COMPREPLY[@]} -eq 0 && ${#noun_aliases[@]} -gt 0 && ${#must_have_one_noun[@]} -ne 0 ]]; then
        COMPREPLY=( $(compgen -W "${noun_aliases[*]}" -- "$cur") )
    fi

    if [[ ${#COMPREPLY[@]} -eq 0 ]]; then
        declare -F __custom_func >/dev/null && __custom_func
    fi

    # available in bash-completion >= 2, not always present on macOS
    if declare -F __ltrim_colon_completions >/dev/null; then
        __ltrim_colon_completions "$cur"
    fi
}

# The arguments should be in the form "ext1|ext2|extn"
__handle_filename_extension_flag()
{
    local ext="$1"
    _filedir "@(${ext})"
}

__handle_subdirs_in_dir_flag()
{
    local dir="$1"
    pushd "${dir}" >/dev/null 2>&1 && _filedir -d && popd >/dev/null 2>&1
}

__handle_flag()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    # if a command required a flag, and we found it, unset must_have_one_flag()
    local flagname=${words[c]}
    local flagvalue
    # if the word contained an =
    if [[ ${words[c]} == *"="* ]]; then
        flagvalue=${flagname#*=} # take in as flagvalue after the =
        flagname=${flagname%%=*} # strip everything after the =
        flagname="${flagname}=" # but put the = back
    fi
    __debug "${FUNCNAME[0]}: looking for ${flagname}"
    if __contains_word "${flagname}" "${must_have_one_flag[@]}"; then
        must_have_one_flag=()
    fi

    # if you set a flag which only applies to this command, don't show subcommands
    if __contains_word "${flagname}" "${local_nonpersistent_flags[@]}"; then
      commands=()
    fi

    # keep flag value with flagname as flaghash
    if [ -n "${flagvalue}" ] ; then
        flaghash[${flagname}]=${flagvalue}
    elif [ -n "${words[ $((c+1)) ]}" ] ; then
        flaghash[${flagname}]=${words[ $((c+1)) ]}
    else
        flaghash[${flagname}]="true" # pad "true" for bool flag
    fi

    # skip the argument to a two word flag
    if __contains_word "${words[c]}" "${two_word_flags[@]}"; then
        c=$((c+1))
        # if we are looking for a flags value, don't show commands
        if [[ $c -eq $cword ]]; then
            commands=()
        fi
    fi

    c=$((c+1))

}

__handle_noun()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    if __contains_word "${words[c]}" "${must_have_one_noun[@]}"; then
        must_have_one_noun=()
    elif __contains_word "${words[c]}" "${noun_aliases[@]}"; then
        must_have_one_noun=()
    fi

    nouns+=("${words[c]}")
    c=$((c+1))
}

__handle_command()
{
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"

    local next_command
    if [[ -n ${last_command} ]]; then
        next_command="_${last_command}_${words[c]//:/__}"
    else
        if [[ $c -eq 0 ]]; then
            next_command="_$(basename "${words[c]//:/__}")"
        else
            next_command="_${words[c]//:/__}"
        fi
    fi
    c=$((c+1))
    __debug "${FUNCNAME[0]}: looking for ${next_command}"
    declare -F "$next_command" >/dev/null && $next_command
}

__handle_word()
{
    if [[ $c -ge $cword ]]; then
        __handle_reply
        return
    fi
    __debug "${FUNCNAME[0]}: c is $c words[c] is ${words[c]}"
    if [[ "${words[c]}" == -* ]]; then
        __handle_flag
    elif __contains_word "${words[c]}" "${commands[@]}"; then
        __handle_command
    elif [[ $c -eq 0 ]] && __contains_word "$(basename "${words[c]}")" "${commands[@]}"; then
        __handle_command
    else
        __handle_noun
    fi
    __handle_word
}

_vpc_agent()
{
    last_command="vpc_agent"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_db_migrate()
{
    last_command="vpc_db_migrate"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_db_ping()
{
    last_command="vpc_db_ping"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_db()
{
    last_command="vpc_db"
    commands=()
    commands+=("migrate")
    commands+=("ping")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_doc_man()
{
    last_command="vpc_doc_man"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--man-dir=")
    two_word_flags+=("-m")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_doc_md()
{
    last_command="vpc_doc_md"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dir=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--dir=")
    flags+=("--url-prefix=")
    local_nonpersistent_flags+=("--url-prefix=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_doc()
{
    last_command="vpc_doc"
    commands=()
    commands+=("man")
    commands+=("md")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_ethlink_destroy()
{
    last_command="vpc_ethlink_destroy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ethlink-id=")
    two_word_flags+=("-E")
    local_nonpersistent_flags+=("--ethlink-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--ethlink-id=")
    must_have_one_flag+=("-E")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_ethlink_list()
{
    last_command="vpc_ethlink_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--sort-by=")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--sort-by=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_ethlink_vtag()
{
    last_command="vpc_ethlink_vtag"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ethlink-id=")
    two_word_flags+=("-E")
    local_nonpersistent_flags+=("--ethlink-id=")
    flags+=("--get-vtag")
    flags+=("-g")
    local_nonpersistent_flags+=("--get-vtag")
    flags+=("--set-vtag=")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--set-vtag=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--ethlink-id=")
    must_have_one_flag+=("-E")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_ethlink()
{
    last_command="vpc_ethlink"
    commands=()
    commands+=("destroy")
    commands+=("list")
    commands+=("vtag")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_hostlink_create()
{
    last_command="vpc_hostlink_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--hostlink-id=")
    two_word_flags+=("-H")
    local_nonpersistent_flags+=("--hostlink-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_hostlink_destroy()
{
    last_command="vpc_hostlink_destroy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--hostlink-id=")
    two_word_flags+=("-H")
    local_nonpersistent_flags+=("--hostlink-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--hostlink-id=")
    must_have_one_flag+=("-H")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_hostlink_genmac()
{
    last_command="vpc_hostlink_genmac"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_hostlink_list()
{
    last_command="vpc_hostlink_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_hostlink()
{
    last_command="vpc_hostlink"
    commands=()
    commands+=("create")
    commands+=("destroy")
    commands+=("genmac")
    commands+=("list")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_interface_list()
{
    last_command="vpc_interface_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_interface()
{
    last_command="vpc_interface"
    commands=()
    commands+=("list")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_list()
{
    last_command="vpc_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--obj-counts")
    flags+=("-c")
    local_nonpersistent_flags+=("--obj-counts")
    flags+=("--obj-type=")
    two_word_flags+=("-t")
    local_nonpersistent_flags+=("--obj-type=")
    flags+=("--sort-by=")
    two_word_flags+=("-s")
    local_nonpersistent_flags+=("--sort-by=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_shell_autocomplete_bash()
{
    last_command="vpc_shell_autocomplete_bash"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--dir=")
    two_word_flags+=("-d")
    local_nonpersistent_flags+=("--dir=")
    flags+=("--help")
    flags+=("-h")
    local_nonpersistent_flags+=("--help")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_shell_autocomplete()
{
    last_command="vpc_shell_autocomplete"
    commands=()
    commands+=("bash")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_shell()
{
    last_command="vpc_shell"
    commands=()
    commands+=("autocomplete")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_create()
{
    last_command="vpc_switch_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--switch-id=")
    local_nonpersistent_flags+=("--switch-id=")
    flags+=("--vni=")
    local_nonpersistent_flags+=("--vni=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--vni=")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_destroy()
{
    last_command="vpc_switch_destroy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--switch-id=")
    local_nonpersistent_flags+=("--switch-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--switch-id=")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_list()
{
    last_command="vpc_switch_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_port_add()
{
    last_command="vpc_switch_port_add"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--ethlink-id=")
    local_nonpersistent_flags+=("--ethlink-id=")
    flags+=("--l2-name=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--l2-name=")
    flags+=("--port-id=")
    local_nonpersistent_flags+=("--port-id=")
    flags+=("--switch-id=")
    local_nonpersistent_flags+=("--switch-id=")
    flags+=("--uplink")
    flags+=("-u")
    local_nonpersistent_flags+=("--uplink")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_port_connect()
{
    last_command="vpc_switch_port_connect"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--interface-id=")
    two_word_flags+=("-I")
    local_nonpersistent_flags+=("--interface-id=")
    flags+=("--port-id=")
    local_nonpersistent_flags+=("--port-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--interface-id=")
    must_have_one_flag+=("-I")
    must_have_one_flag+=("--port-id=")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_port_disconnect()
{
    last_command="vpc_switch_port_disconnect"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--interface-id=")
    two_word_flags+=("-I")
    local_nonpersistent_flags+=("--interface-id=")
    flags+=("--port-id=")
    local_nonpersistent_flags+=("--port-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--interface-id=")
    must_have_one_flag+=("-I")
    must_have_one_flag+=("--port-id=")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_port_remove()
{
    last_command="vpc_switch_port_remove"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--port-id=")
    local_nonpersistent_flags+=("--port-id=")
    flags+=("--switch-id=")
    local_nonpersistent_flags+=("--switch-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--port-id=")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch_port()
{
    last_command="vpc_switch_port"
    commands=()
    commands+=("add")
    commands+=("connect")
    commands+=("disconnect")
    commands+=("remove")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_switch()
{
    last_command="vpc_switch"
    commands=()
    commands+=("create")
    commands+=("destroy")
    commands+=("list")
    commands+=("port")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_version()
{
    last_command="vpc_version"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vm_create()
{
    last_command="vpc_vm_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vm()
{
    last_command="vpc_vm"
    commands=()
    commands+=("create")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_create()
{
    last_command="vpc_vmnic_create"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--vmnic-id=")
    two_word_flags+=("-N")
    local_nonpersistent_flags+=("--vmnic-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_destroy()
{
    last_command="vpc_vmnic_destroy"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--vmnic-id=")
    two_word_flags+=("-N")
    local_nonpersistent_flags+=("--vmnic-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--vmnic-id=")
    must_have_one_flag+=("-N")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_genmac()
{
    last_command="vpc_vmnic_genmac"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_get()
{
    last_command="vpc_vmnic_get"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--num-queues")
    flags+=("-n")
    local_nonpersistent_flags+=("--num-queues")
    flags+=("--vmnic-id=")
    two_word_flags+=("-N")
    local_nonpersistent_flags+=("--vmnic-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--vmnic-id=")
    must_have_one_flag+=("-N")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_list()
{
    last_command="vpc_vmnic_list"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic_set()
{
    last_command="vpc_vmnic_set"
    commands=()

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--freeze")
    flags+=("-E")
    local_nonpersistent_flags+=("--freeze")
    flags+=("--num-queues=")
    two_word_flags+=("-n")
    local_nonpersistent_flags+=("--num-queues=")
    flags+=("--unfreeze")
    local_nonpersistent_flags+=("--unfreeze")
    flags+=("--vmnic-id=")
    two_word_flags+=("-N")
    local_nonpersistent_flags+=("--vmnic-id=")
    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_flag+=("--vmnic-id=")
    must_have_one_flag+=("-N")
    must_have_one_noun=()
    noun_aliases=()
}

_vpc_vmnic()
{
    last_command="vpc_vmnic"
    commands=()
    commands+=("create")
    commands+=("destroy")
    commands+=("genmac")
    commands+=("get")
    commands+=("list")
    commands+=("set")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

_vpc()
{
    last_command="vpc"
    commands=()
    commands+=("agent")
    commands+=("db")
    commands+=("doc")
    commands+=("ethlink")
    commands+=("hostlink")
    commands+=("interface")
    commands+=("list")
    commands+=("shell")
    commands+=("switch")
    commands+=("version")
    commands+=("vm")
    commands+=("vmnic")

    flags=()
    two_word_flags=()
    local_nonpersistent_flags=()
    flags_with_completion=()
    flags_completion=()

    flags+=("--log-format=")
    two_word_flags+=("-F")
    flags+=("--log-level=")
    two_word_flags+=("-l")
    flags+=("--use-color")
    flags+=("--use-pager")
    flags+=("-P")
    flags+=("--utc")
    flags+=("-Z")

    must_have_one_flag=()
    must_have_one_noun=()
    noun_aliases=()
}

__start_vpc()
{
    local cur prev words cword
    declare -A flaghash 2>/dev/null || :
    if declare -F _init_completion >/dev/null 2>&1; then
        _init_completion -s || return
    else
        __my_init_completion -n "=" || return
    fi

    local c=0
    local flags=()
    local two_word_flags=()
    local local_nonpersistent_flags=()
    local flags_with_completion=()
    local flags_completion=()
    local commands=("vpc")
    local must_have_one_flag=()
    local must_have_one_noun=()
    local last_command
    local nouns=()

    __handle_word
}

if [[ $(type -t compopt) = "builtin" ]]; then
    complete -o default -F __start_vpc vpc
else
    complete -o default -o nospace -F __start_vpc vpc
fi

# ex: ts=4 sw=4 et filetype=sh
