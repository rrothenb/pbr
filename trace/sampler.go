package trace

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const threshold = 0.25
const adapt = 2
const block = 5

// Sampler traces samples from light paths in a scene
type Sampler struct {
	Width   int
	Height  int
	pixels  []float64 // r, g, b, count
	cam     *Camera
	scene   *Scene
	bounces int
	count   int
}

// NewSampler constructs a new Sampler instance
func NewSampler(cam *Camera, scene *Scene, bounces int) *Sampler {
	return &Sampler{
		Width:   cam.Width,
		Height:  cam.Height,
		pixels:  make([]float64, cam.Width*cam.Height*block),
		cam:     cam,
		scene:   scene,
		bounces: bounces,
	}
}

// Collect traces light paths for the full image
func (s *Sampler) Collect(frames int, samples int) {
	results := make(chan []float64)
	for i := 0; i < frames; i++ {
		go s.scan(samples, results)
	}
	for i := 0; i < frames; i++ {
		result := <-results
		fmt.Printf("Frame %v/%v complete.\n", i, frames)
		for p := 0; p < len(result); p++ {
			s.pixels[p] += result[p]
		}
	}
}

// Scan takes samples of every pixel in the image
func (s *Sampler) scan(samples int, result chan []float64) {
	pixels := make([]float64, s.Width*s.Height*block)
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	sampled := 0
	for p := 0; p < len(pixels); p += block {
		s.sample(pixels, p, rnd, 1)
		sampled++
	}
	delta := 0.0
	for p := 0; p < len(pixels); p += block {
		delta += s.sample(pixels, p, rnd, 1)
		sampled++
	}
	for sampled < samples {
		size := float64(len(pixels) / block)
		mean := delta / size
		fmt.Println(sampled, samples, float64(sampled)/float64(samples), delta, size, mean)
		delta = 0.0
		for p := 0; sampled < samples && p < len(pixels); p += block {
			adapted := int(math.Ceil(pixels[p+4] / mean))
			delta += s.sample(pixels, p, rnd, adapted)
			sampled += adapted
		}
	}
	result <- pixels
}

func (s *Sampler) sample(pixels []float64, p int, rnd *rand.Rand, samples int) float64 {
	x, y := s.offsetPixel(p)
	before := value(pixels, p)
	for i := 0; i < samples; i++ {
		sample := s.trace(x, y, rnd)
		rgb := sample.Array()
		pixels[p] += rgb[0]
		pixels[p+1] += rgb[1]
		pixels[p+2] += rgb[2]
		pixels[p+3]++
	}
	after := value(pixels, p)
	delta := before.Minus(after).Length()
	pixels[p+4] = delta
	return pixels[p+4]
}

func value(pixels []float64, i int) Vector3 {
	if pixels[i+3] == 0 {
		return Vector3{}
	}
	sample := Vector3{pixels[i], pixels[i+1], pixels[i+2]}
	return sample.Scale(1 / pixels[i+3])
}

func (s *Sampler) trace(x, y int, rnd *rand.Rand) Vector3 {
	ray := s.cam.Ray(x, y, rnd)
	signal := Vector3{1, 1, 1}
	energy := Vector3{0, 0, 0}

	for bounce := 0; bounce < s.bounces; bounce++ {
		intersected, hit := s.scene.Intersect(ray)
		if !intersected {
			energy = energy.Add(s.scene.Env(ray).Mult(signal))
			break
		}
		light := hit.Mat.Emit(hit.Normal, ray.Dir)
		energy = energy.Add(light.Mult(signal))
		if rnd.Float64() > signal.Max() {
			break
		}
		signal = signal.Scale(1 / signal.Max())
		next, dir, strength := hit.Mat.Bsdf(hit.Normal, ray.Dir, hit.Dist, rnd)
		if !next {
			break
		}
		ray = Ray3{hit.Point, dir}
		signal = signal.Mult(strength)
	}

	return energy
}

func (s *Sampler) offsetPixel(i int) (x, y int) {
	pos := i / block
	return pos % s.Width, pos / s.Width
}

// Values gets the average sampled rgb at each pixel
func (s *Sampler) Values() []float64 {
	rgb := make([]float64, s.Width*s.Height*3)
	for p := 0; p < len(s.pixels); p += block {
		val := value(s.pixels, p).Array()
		i := p / block * 3
		rgb[i] = val[0]
		rgb[i+1] = val[1]
		rgb[i+2] = val[2]
	}
	return rgb
}

// Counts returns the sample count at each pixel as rgb
func (s *Sampler) Counts() []float64 {
	rgb := make([]float64, s.Width*s.Height*3)
	var max float64
	for p := 0; p < len(s.pixels); p += block {
		max = math.Max(max, s.pixels[p+3])
	}
	for p := 0; p < len(s.pixels); p += block {
		val := (s.pixels[p+3] / max) * 255
		i := p / block * 3
		rgb[i] = val
		rgb[i+1] = val
		rgb[i+2] = val
	}
	return rgb
}
