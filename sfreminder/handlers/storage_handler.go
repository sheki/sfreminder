//  handlerStorage used for an non related experiment. 
func handleStorage(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	e1 := Employee{
		Name:     "Suppa James",
		Role:     "Manager",
		HireDate: time.Now(),
	}

	_, err := datastore.Put(c, datastore.NewIncompleteKey(c, "employee", nil), &e1)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "Stored and retrieved the Employee named %q", e1.Name)
}