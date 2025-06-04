let split_comma s =
  String.split_on_char ',' s |> List.map String.trim |> List.filter (fun x -> x <> "")
