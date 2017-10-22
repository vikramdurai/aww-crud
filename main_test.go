package main

import "testing"

func TestSlug(t *testing.T) {
	tt := []struct {
		name  string
		title string
		slug  string
	}{
		{"&, :, and ()", "Jack & Jill: The Untold Story (Part 2)", "jack--jill-the-untold-story-part-2"},
		{"! and @", "Email me at me@myself.com", "email-me-at-memyself.com"},
		{"$ and %", "$50 is 50% off $100", "50-is-50-off-100"},
		{"kitchen sink", "Raindrops & Roses: Hey! me@myself.com #1 $300.00 10% off x^2 7 * 6 = 42 (the untold story)", "raindrops--roses-hey-memyself.com-1-300.00-10-off-x2-7--6-=-42-the-untold-story"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			rec := &Record{Title: tc.title, Content: ""}
			s := rec.Slug();
			if s != tc.slug {
				t.Fatalf("\nexpected: %s\nactual: %s", tc.slug, s)
			}
		})
	}
}