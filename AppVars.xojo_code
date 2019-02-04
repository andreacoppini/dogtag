#tag Module
Protected Module AppVars
	#tag Method, Flags = &h0
		Function NowDate() As String
		  Dim v As Variant
		  Dim d As Date
		  
		  d = New Date
		  Return d.ShortDate
		  
		  
		End Function
	#tag EndMethod

	#tag Method, Flags = &h0
		Function NowTime() As String
		  Dim v As Variant
		  Dim d As Date
		  
		  d = New Date
		  Return d.ShortTime
		  
		  
		End Function
	#tag EndMethod

	#tag Method, Flags = &h0
		Sub PopSelByText(Extends p As PopupMenu, text As String)
		  For i As Integer = 0 To p.ListCount-1
		    If p.List(i) = text Then
		      p.ListIndex = i
		      Return
		    End If
		  Next
		End Sub
	#tag EndMethod


	#tag Property, Flags = &h0
		APCSVList(-1,-1) As String
	#tag EndProperty

	#tag Property, Flags = &h0
		APGroupArr(-1,-1) As String
	#tag EndProperty

	#tag Property, Flags = &h0
		APGroupSel As String
	#tag EndProperty

	#tag Property, Flags = &h0
		arrorideURL(-1) As String
	#tag EndProperty

	#tag Property, Flags = &h0
		arrorideURLList(-1) As String
	#tag EndProperty

	#tag Property, Flags = &h0
		Green As Color = &c00ff00
	#tag EndProperty

	#tag Property, Flags = &h0
		Red As Color = &cff0000
	#tag EndProperty

	#tag Property, Flags = &h0
		SZCookie As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZSessionAll As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZURL As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZVer As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZZoneArr(-1,-1) As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZZones As String
	#tag EndProperty

	#tag Property, Flags = &h0
		SZZoneSel As String
	#tag EndProperty


	#tag ViewBehavior
		#tag ViewProperty
			Name="Name"
			Visible=true
			Group="ID"
			Type="String"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Index"
			Visible=true
			Group="ID"
			InitialValue="-2147483648"
			Type="Integer"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Super"
			Visible=true
			Group="ID"
			Type="String"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Left"
			Visible=true
			Group="Position"
			InitialValue="0"
			Type="Integer"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Top"
			Visible=true
			Group="Position"
			InitialValue="0"
			Type="Integer"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZURL"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZZones"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZZoneSel"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="APGroupSel"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZVer"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZCookie"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="SZSessionAll"
			Group="Behavior"
			Type="String"
			EditorType="MultiLineEditor"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Red"
			Group="Behavior"
			InitialValue="&cff0000"
			Type="Color"
		#tag EndViewProperty
		#tag ViewProperty
			Name="Green"
			Group="Behavior"
			InitialValue="&c00ff00"
			Type="Color"
		#tag EndViewProperty
	#tag EndViewBehavior
End Module
#tag EndModule
