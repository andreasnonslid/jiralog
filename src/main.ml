open Lwt.Infix

let jira_host () =
  match Sys.getenv_opt "JIRA_HOST" with
  | Some h -> h
  | None -> "autostore.atlassian.net"

let get_auth_header () =
  match Sys.getenv_opt "JIRA_EMAIL", Sys.getenv_opt "JIRA_TOKEN" with
  | Some email, Some token ->
      let creds = email ^ ":" ^ token in
      "Basic " ^ Base64.encode_string creds
  | _ -> failwith "JIRA_EMAIL or JIRA_TOKEN not set"

let jira_request ~meth ~path ~body =
  let uri = Uri.make ~scheme:"https" ~host:(jira_host ()) ~path in
  get_auth_header () >>= fun auth ->
  let headers = Cohttp.Header.init_with "Authorization" auth in
  let headers = Cohttp.Header.add headers "Accept" "application/json" in
  let headers = Cohttp.Header.add headers "Content-Type" "application/json" in
  Cohttp_lwt_unix.Client.call ~headers ~body:(`String body) meth uri
  >>= fun (_, body) -> Cohttp_lwt.Body.to_string body

let search jql fields =
  let payload = `Assoc [
      ("jql", `String jql);
      ("fields", `List (List.map (fun f -> `String f) fields))
    ] |> Yojson.Safe.to_string in
  jira_request ~meth:`POST ~path:"/rest/api/3/search" ~body:payload
  >|= fun data ->
  let json = Yojson.Safe.from_string data in
  let issues = json
    |> Yojson.Safe.Util.member "issues"
    |> Yojson.Safe.Util.to_list in
  List.iter (fun issue ->
    let key = Yojson.Safe.Util.member "key" issue |> Yojson.Safe.Util.to_string in
    let summary =
      issue |> Yojson.Safe.Util.member "fields"
      |> Yojson.Safe.Util.member "summary" |> Yojson.Safe.Util.to_string in
    Printf.printf "%s : %s\n" key summary
  ) issues

let log_work issue started seconds comment =
  let body = `Assoc [
    ("comment", `Assoc [
      ("type", `String "doc");
      ("version", `Int 1);
      ("content", `List [
        `Assoc [
          ("type", `String "paragraph");
          ("content", `List [ `Assoc [ ("type", `String "text"); ("text", `String comment) ] ])
        ]
      ])
    ]);
    ("started", `String started);
    ("timeSpentSeconds", `Int seconds);
    ("adjustEstimate", `String "auto")
  ] |> Yojson.Safe.to_string in
  let path = Printf.sprintf "/rest/api/3/issue/%s/worklog?notifyUsers=false" issue in
  jira_request ~meth:`POST ~path ~body
  >|= fun _ -> ()

open Cmdliner

let jql =
  let doc = "JQL query" in
  Arg.(required & opt (some string) None & info ["jql"] ~doc)

let fields =
  let doc = "Comma separated fields" in
  Arg.(value & opt string "summary" & info ["fields"] ~doc)

let search_cmd =
  let run jql fields =
    let fields = String.split_on_char ',' fields |> List.map String.trim in
    Lwt_main.run (search jql fields)
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
    Lwt_main.run (log_work issue started seconds comment)
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
