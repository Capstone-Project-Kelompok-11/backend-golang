package util

type Rating struct {
  Rating1 int
  Rating2 int
  Rating3 int
  Rating4 int
  Rating5 int
}

func RatingView(rating Rating) float64 {

  var z float64

  r1 := float64(1 * rating.Rating1)
  r2 := float64(2 * rating.Rating2)
  r3 := float64(3 * rating.Rating3)
  r4 := float64(4 * rating.Rating4)
  r5 := float64(5 * rating.Rating5)
  n := float64(rating.Rating1 + rating.Rating2 + rating.Rating3 + rating.Rating4 + rating.Rating5)
  k := (r1 + r2 + r3 + r4 + r5) / n
  x := []float64{.5, 1, 1.5, 2, 2.5, 3, 3.5, 4, 4.5, 5}

  z = 0
  for i := 0; i < len(x); i++ {

    if x[i] <= k {

      z = x[i]
      continue
    }

    return z
  }

  return z
}
