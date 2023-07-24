open System
open System.IO

let parse lines =
    seq {
        let rows = lines |> Seq.mapi (fun y row -> (y, row))

        for (y, row) in rows do
            let cols = row |> Seq.mapi (fun x col -> (x, col))

            for (x, col) in cols do
                if col = '#' then
                    (x, y)
    }
    |> Set.ofSeq

let between (xa, ya) (xb, yb) (xc, yc) =
    let dxc, dyc = xc - xa, yc - ya
    let dxl, dyl = xb - xa, yb - ya
    let inLineOfSight = (dxc * dyl - dyc * dxl) = 0

    let xmin, xmax = (min xa xb), (max xa xb)
    let ymin, ymax = (min ya yb), (max ya yb)
    let betweenPoints = xc >= xmin && xc <= xmax && yc >= ymin && yc <= ymax

    inLineOfSight && betweenPoints


let asteroids = File.ReadAllLines "in" |> parse


let ans1 =
    asteroids
    |> Seq.map (fun a ->
        a,
        asteroids - Set([ a ])
        |> Seq.map (fun b ->
            asteroids - Set([ a; b ])
            |> Seq.forall (not << between a b)
            |> fun canDetect ->
                match canDetect with
                | true -> Some(b)
                | false -> None)
        |> Seq.choose id
        |> Seq.length)
    |> Seq.maxBy snd

printfn $"ans1 {ans1}"


let angle (x, y) =
    // unit vector pointing towards -y
    let ux, uy = 0, -1

    // calculate angle between -y axis and asteroid
    // https://stackoverflow.com/questions/14066933/direct-way-of-computing-the-clockwise-angle-between-two-vectors#16544330
    let dot = ux * x + uy * y
    let det = ux * y - uy * x
    let rad = Math.Atan2(float det, float dot)
    let deg = rad * (180.0 / Math.PI)

    match deg with
    | x when x < 0 -> x + 360.0
    | _ -> deg

let distance (x, y) = Math.Sqrt(float x ** 2 + float y ** 2)

let ans2 =
    // move station to center (0,0) and asteroids relative to it
    let sx, sy = fst ans1
    let asteroids = asteroids |> Set.map (fun (x, y) -> x - sx, y - sy)

    let mutable targets = asteroids - Set([ (0, 0) ])
    let mutable destroyed = [ (0, 0) ]

    while Seq.length targets > 0 do
        let mutable angles = targets |> Seq.groupBy (angle) |> Seq.sortBy fst |> Seq.toList

        for _, canidates in angles do
            let destroy = canidates |> Seq.minBy (distance)
            destroyed <- destroyed @ [ destroy ]

        targets <- targets - Set.ofList destroyed

    let ax, ay = fst destroyed[200] + sx, snd destroyed[200] + sy
    ax * 100 + ay

printfn $"ans2 {ans2}"
