package bubbleforms

//func Test_listModel_Update(t *testing.T) {
//	type fields struct {
//		list list.Model
//	}
//	type args struct {
//		msg tea.Msg
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		args   args
//		want   bool
//		want1  bool
//	}{
//		{
//			name:   "char",
//			fields: fields{},
//			args:   args{msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")}},
//			want:   true,
//			want1:  false,
//		},
//		{
//			name:   "enter",
//			fields: fields{},
//			args:   args{msg: tea.KeyMsg{Type: tea.KeyEnter}},
//			want:   true,
//			want1:  false,
//		},
//		{
//			name:   "quit",
//			fields: fields{},
//			args:   args{msg: tea.KeyMsg{Type: tea.KeyCtrlC}},
//			want:   true,
//			want1:  true,
//		},
//		{
//			name:   "window size",
//			fields: fields{},
//			args:   args{msg: tea.WindowSizeMsg{Width: 80, Height: 20}},
//			want:   true,
//			want1:  false,
//		},
//	}
//	m := NewListModel()
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, got1 := m.Update(tt.args.msg)
//			if tt.want != (got != nil) {
//				t.Errorf("Update() got = %v, want %v", got, tt.want)
//			}
//			if tt.want1 != (got1 != nil) {
//				t.Errorf("Update() got1 = %v, want1 %v", got1, tt.want1)
//			}
//		})
//	}
//}
//
//func Test_newListModel(t *testing.T) {
//	m := NewListModel()
//	m.Init()
//	i, _ := m.List.SelectedItem().(ListItem)
//	if i.Title() == "" {
//		t.Errorf("empty item Title")
//	}
//	if i.Description() == "" {
//		t.Errorf("empty item Description")
//	}
//	if i.FilterValue() == "" {
//		t.Errorf("empty item FilterValue")
//	}
//	if m.View() == "" {
//		t.Errorf("empty model View")
//	}
//}
