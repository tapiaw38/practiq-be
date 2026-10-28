package domain

func StudentsToDeactivate(byActivity []string, maxStudents int, chosen []string) []string {
	if maxStudents < 0 {
		maxStudents = 0
	}
	if len(byActivity) <= maxStudents {
		return nil
	}

	if len(chosen) > 0 {
		keep := make(map[string]bool, len(chosen))

		for i, id := range chosen {
			if i >= maxStudents {
				break
			}
			keep[id] = true
		}

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
