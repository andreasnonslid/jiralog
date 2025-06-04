
open Cmdliner
open Jiralog

let jql =
  let doc = "JQL query" in
  Arg.(required & opt (some string) None & info ["jql"] ~doc)

let fields =
  let doc = "Comma separated fields" in
  Arg.(value & opt string "summary" & info ["fields"] ~doc)

let search_cmd =
  let run jql fields =
    let fields = Util.split_comma fields in
    Lwt_main.run (Jira.search jql fields)
  in
  Term.(const run $ jql $ fields),
  Term.info "search" ~doc:"Search issues"

let issue =
  let doc = "Issue key" in
  Arg.(required & opt (some string) None & info ["issue"] ~doc)

let started =
  let doc = "Start time RFC3339" in
  Arg.(required & opt (some string) None & info ["started"] ~doc)

let seconds =
  let doc = "Time spent in seconds" in
  Arg.(required & opt (some int) None & info ["seconds"] ~doc)

let comment =
  let doc = "Comment text" in
  Arg.(value & opt string "" & info ["comment"] ~doc)

let log_cmd =
  let run issue started seconds comment =
    Lwt_main.run (Jira.log_work issue started seconds comment)
  in
  Term.(const run $ issue $ started $ seconds $ comment),
  Term.info "log" ~doc:"Log work"

let () =
  let cmds = [search_cmd; log_cmd] in
  let default =
    let doc = "Simple CLI to interact with Jira" in
    Term.(ret (const (`Help (`Pager, None)))), Term.info "jiralog" ~version:"0.1" ~doc
  in
  match Term.eval_choice default cmds with
  | `Error _ -> exit 1
  | _ -> exit 0
