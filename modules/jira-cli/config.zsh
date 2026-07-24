#!/bin/bash
#
alias jira_me_assigned="jira issue list -a$(jira me)"
alias jira_me_reported="jira issue list -r$(jira me)"
alias jira_me_reported_week="jira issue list -r$(jira me) --created week"

alias jira_issue_history="jira issue list --history"

alias jira_sprint_active="jira sprint list --current"
alias jira_sprint_me="jira sprint list --current -a$(jira me)"
alias jira_sprint_next="jira sprint list --next"
alias jira_sprint_prev="jira sprint list --prev"

alias jira_board_list="jira board list"
