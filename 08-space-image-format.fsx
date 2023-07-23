open System
open System.IO

let log x = printfn "%A" x

let w = 25
let h = 6

let size = w * h

let input =
    File.ReadAllLines "in"
    |> Seq.head
    |> Seq.map Char.GetNumericValue
    |> Seq.map int
    |> Seq.toList

let chunk list size =
    list
    |> Seq.mapi (fun i x -> i / size, x)
    |> Seq.groupBy fst
    |> Seq.map (fun (_, g) -> Seq.map snd g |> Seq.toList)
    |> Seq.toList

let layers = chunk input size

let countDigit list digit =
    list |> Seq.filter (fun x -> x = digit) |> Seq.length

let layerFewest0 = layers |> Seq.minBy (fun layer -> countDigit layer 0)

let ans1 = (countDigit layerFewest0 1) * (countDigit layerFewest0 2)
log ans1

let calcPixel pxim pxla =
    match pxim with
    | 2 -> pxla
    | _ -> pxim

let readablePixel px =
    match px with
    | 1 -> "#"
    | _ -> " "

let image =
    List.fold (fun image layer -> layer |> List.mapi (fun i p -> calcPixel image[i] layer[i])) layers[0] layers[1..]

chunk image w
|> Seq.map (fun row -> row |> Seq.map readablePixel |> Seq.fold (+) "")
|> Seq.iter log
