package domain

// StudentsToDeactivate picks who stops counting when a plan gets smaller.
//
// `byActivity` is the school's active students, least recently active first.
// That order is the point: a teacher's oldest students are usually the ones who
// already finished, and their newest are the class sitting in front of them.
// Cutting by seniority would keep the finished ones and deactivate the working
// ones, which is backwards.
//
// The teacher can override the choice. This is only what happens if they do
// not, and it runs when the period they already paid for ends — cutting the
// day they downgrade would take away something they bought, and would punish
// students who did not make the decision.
func StudentsToDeactivate(byActivity []string, maxStudents int, chosen []string) []string {
	if maxStudents < 0 {
		maxStudents = 0
	}
	if len(byActivity) <= maxStudents {
		return nil
	}

	if len(chosen) > 0 {
		keep := make(map[string]bool, len(chosen))
		// More chosen than the plan allows means the choice cannot be honoured
		// as given; the extra ones fall back to the automatic order below.
		for i, id := range chosen {
			if i >= maxStudents {
				break
			}
			keep[id] = true
		}

		// Any room the choice left over is filled from the most recently
		// active end, the same rule the automatic order follows. Filling from
		// the other end would keep the dormant students the rule exists to cut.
		room := maxStudents - len(keep)
		for i := len(byActivity) - 1; i >= 0 && room > 0; i-- {
			if keep[byActivity[i]] {
				continue
			}
			keep[byActivity[i]] = true
			room--
		}

		out := make([]string, 0, len(byActivity)-len(keep))
		for _, id := range byActivity {
			if !keep[id] {
				out = append(out, id)
			}
		}
		return out
	}

	return append([]string(nil), byActivity[:len(byActivity)-maxStudents]...)
}
