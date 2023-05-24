package lik

func CompareItems(itema Itemer, itemb Itemer) bool {
	if itema == nil && itemb == nil {
		return true
	} else if itema == nil || itemb == nil {
		return false
	} else if itema.IsSet() && itemb.IsSet() {
		return CompareSets(itema.ToSet(), itemb.ToSet())
	} else if itema.IsSet() || itemb.IsSet() {
		return false
	} else if itema.IsList() && itemb.IsList() {
		return CompareLists(itema.ToList(), itemb.ToList())
	} else if itema.IsList() || itemb.IsList() {
		return false
	} else {
		return itema.ToString() == itemb.ToString()
	}
}

func CompareSets(seta Seter, setb Seter) bool {
	if seta == nil && setb == nil {
		return true
	} else if seta == nil || setb == nil {
		return false
	} else if keysa, keysb := seta.SortKeys(), setb.SortKeys(); len(keysa) != len(keysb) {
		return false
	} else {
		for k, key := range keysa {
			if keyb := keysb[k]; keyb != key {
				return false
			} else if vala, valb := seta.GetItem(key), setb.GetItem(keyb); !CompareItems(vala, valb) {
				return false
			}
		}
	}
	return true
}

func CompareLists(lista Lister, listb Lister) bool {
	if lista == nil && listb == nil {
		return true
	} else if lista == nil || listb == nil {
		return false
	} else if lena, lenb := lista.Count(), listb.Count(); lena != lenb {
		return false
	} else {
		for k := 0; k < lena; k++ {
			if vala, valb := lista.GetItem(k), listb.GetItem(k); !CompareItems(vala, valb) {
				return false
			}
		}
	}
	return true
}
