open Alcotest
open Jiralog.Util

let test_split () =
  let res = split_comma "a, b, c" in
  check (list string) "fields" ["a"; "b"; "c"] res

let () =
  run "util" [ ("split", [ test_case "basic" `Quick test_split ]) ]
