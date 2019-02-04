#tag Window
Begin Window DogTag
   BackColor       =   &cFFFFFF00
   Backdrop        =   0
   CloseButton     =   True
   Compatibility   =   ""
   Composite       =   False
   Frame           =   0
   FullScreen      =   False
   FullScreenButton=   False
   HasBackColor    =   False
   Height          =   500
   ImplicitInstance=   True
   LiveResize      =   True
   MacProcID       =   0
   MaxHeight       =   32000
   MaximizeButton  =   True
   MaxWidth        =   32000
   MenuBar         =   1855635455
   MenuBarVisible  =   True
   MinHeight       =   500
   MinimizeButton  =   True
   MinWidth        =   650
   Placement       =   0
   Resizeable      =   True
   Title           =   "DogTag 5 for Ruckus SmartZone"
   Visible         =   True
   Width           =   650
   Begin TextArea txtDebug
      AcceptTabs      =   False
      Alignment       =   0
      AutoDeactivate  =   False
      AutomaticallyCheckSpelling=   False
      BackColor       =   &cFFFFFF00
      Bold            =   False
      Border          =   True
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Format          =   ""
      Height          =   71
      HelpTag         =   ""
      HideSelection   =   True
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   20
      LimitText       =   0
      LineHeight      =   0.0
      LineSpacing     =   1.0
      LockBottom      =   True
      LockedInPosition=   False
      LockLeft        =   True
      LockRight       =   True
      LockTop         =   False
      Mask            =   ""
      Multiline       =   True
      ReadOnly        =   False
      Scope           =   2
      ScrollbarHorizontal=   True
      ScrollbarVertical=   True
      Styled          =   True
      TabIndex        =   7
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   ""
      TextColor       =   &c00000000
      TextFont        =   "SmallSystem"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   419
      Transparent     =   False
      Underline       =   False
      UseFocusRing    =   True
      Visible         =   False
      Width           =   610
   End
   Begin HTTPSecureSocket SZ
      CertificateFile =   
      CertificatePassword=   ""
      CertificateRejectionFile=   
      ConnectionType  =   5
      Index           =   -2147483648
      LockedInPosition=   False
      Scope           =   0
      Secure          =   True
      TabPanelIndex   =   0
   End
   Begin TextField txtSZUName
      AcceptTabs      =   False
      Alignment       =   0
      AutoDeactivate  =   True
      AutomaticallyCheckSpelling=   False
      BackColor       =   &cFFFFFF00
      Bold            =   False
      Border          =   True
      CueText         =   "Username"
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Format          =   ""
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      Italic          =   False
      Left            =   306
      LimitText       =   0
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Mask            =   ""
      Password        =   False
      ReadOnly        =   False
      Scope           =   0
      TabIndex        =   2
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   "admin"
      TextColor       =   &c00000000
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   14
      Transparent     =   False
      Underline       =   False
      UseFocusRing    =   True
      Visible         =   True
      Width           =   80
   End
   Begin TextField txtSZPwd
      AcceptTabs      =   False
      Alignment       =   0
      AutoDeactivate  =   True
      AutomaticallyCheckSpelling=   False
      BackColor       =   &cFFFFFF00
      Bold            =   False
      Border          =   True
      CueText         =   "Password"
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Format          =   ""
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      Italic          =   False
      Left            =   398
      LimitText       =   0
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Mask            =   ""
      Password        =   True
      ReadOnly        =   False
      Scope           =   0
      TabIndex        =   3
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   ""
      TextColor       =   &c00000000
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   14
      Transparent     =   False
      Underline       =   False
      UseFocusRing    =   True
      Visible         =   True
      Width           =   80
   End
   Begin PushButton btnConnect
      AutoDeactivate  =   True
      Bold            =   False
      ButtonStyle     =   "0"
      Cancel          =   False
      Caption         =   "Connect"
      Default         =   False
      Enabled         =   True
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   505
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Scope           =   0
      TabIndex        =   5
      TabPanelIndex   =   0
      TabStop         =   True
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   14
      Transparent     =   False
      Underline       =   False
      Visible         =   True
      Width           =   98
   End
   Begin Oval ovlConnect
      AutoDeactivate  =   True
      BorderColor     =   &c00000000
      BorderWidth     =   1
      Enabled         =   True
      FillColor       =   &cFF000000
      Height          =   15
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Left            =   615
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Scope           =   2
      TabIndex        =   -1
      TabPanelIndex   =   0
      Top             =   17
      Transparent     =   False
      Visible         =   True
      Width           =   15
   End
   Begin PopupMenu popZone
      AutoDeactivate  =   False
      Bold            =   False
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   False
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      InitialValue    =   ""
      Italic          =   False
      Left            =   62
      ListIndex       =   0
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   True
      LockRight       =   True
      LockTop         =   True
      Scope           =   2
      TabIndex        =   9
      TabPanelIndex   =   0
      TabStop         =   True
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   48
      Transparent     =   False
      Underline       =   False
      Visible         =   True
      Width           =   324
   End
   Begin Label lblZone
      AutoDeactivate  =   True
      Bold            =   False
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Height          =   20
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   20
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   True
      LockRight       =   False
      LockTop         =   True
      Multiline       =   False
      Scope           =   2
      Selectable      =   False
      TabIndex        =   -1
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   "Zone"
      TextAlign       =   0
      TextColor       =   &c00000000
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   49
      Transparent     =   True
      Underline       =   False
      Visible         =   True
      Width           =   44
   End
   Begin TextField txtSZPort
      AcceptTabs      =   False
      Alignment       =   0
      AutoDeactivate  =   True
      AutomaticallyCheckSpelling=   False
      BackColor       =   &cFFFFFF00
      Bold            =   False
      Border          =   True
      CueText         =   "Port"
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Format          =   ""
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      Italic          =   False
      Left            =   246
      LimitText       =   0
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Mask            =   ""
      Password        =   False
      ReadOnly        =   False
      Scope           =   0
      TabIndex        =   1
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   "8443"
      TextColor       =   &c00000000
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   14
      Transparent     =   False
      Underline       =   False
      UseFocusRing    =   True
      Visible         =   True
      Width           =   48
   End
   Begin CheckBox chkHidePw
      AutoDeactivate  =   True
      Bold            =   False
      Caption         =   ""
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Height          =   20
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   484
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Scope           =   0
      State           =   1
      TabIndex        =   4
      TabPanelIndex   =   0
      TabStop         =   True
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   15
      Transparent     =   False
      Underline       =   False
      Value           =   True
      Visible         =   True
      Width           =   17
   End
   Begin TextField txtSZIP
      AcceptTabs      =   False
      Alignment       =   0
      AutoDeactivate  =   True
      AutomaticallyCheckSpelling=   False
      BackColor       =   &cFFFFFF00
      Bold            =   False
      Border          =   True
      CueText         =   "IP Address"
      DataField       =   ""
      DataSource      =   ""
      Enabled         =   True
      Format          =   ""
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      Italic          =   False
      Left            =   20
      LimitText       =   0
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   True
      LockRight       =   True
      LockTop         =   True
      Mask            =   ""
      Password        =   False
      ReadOnly        =   False
      Scope           =   0
      TabIndex        =   0
      TabPanelIndex   =   0
      TabStop         =   True
      Text            =   ""
      TextColor       =   &c00000000
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   14
      Transparent     =   False
      Underline       =   False
      UseFocusRing    =   True
      Visible         =   True
      Width           =   214
   End
   Begin ProgressBar progBar
      AutoDeactivate  =   True
      Enabled         =   True
      Height          =   21
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Left            =   505
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Maximum         =   100
      Scope           =   0
      TabIndex        =   12
      TabPanelIndex   =   0
      Top             =   48
      Transparent     =   False
      Value           =   0
      Visible         =   False
      Width           =   125
   End
   Begin PushButton btnCSVExport
      AutoDeactivate  =   True
      Bold            =   False
      ButtonStyle     =   "0"
      Cancel          =   False
      Caption         =   "Export CSV..."
      Default         =   False
      Enabled         =   False
      Height          =   22
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   398
      LockBottom      =   False
      LockedInPosition=   False
      LockLeft        =   False
      LockRight       =   True
      LockTop         =   True
      Scope           =   0
      TabIndex        =   13
      TabPanelIndex   =   0
      TabStop         =   True
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   48
      Transparent     =   False
      Underline       =   False
      Visible         =   True
      Width           =   95
   End
   Begin Timer timeKeepalive
      Index           =   -2147483648
      LockedInPosition=   False
      Mode            =   2
      Period          =   10000
      Scope           =   0
      TabPanelIndex   =   0
   End
   Begin TabPanel tabMain
      AutoDeactivate  =   True
      Bold            =   False
      Enabled         =   True
      Height          =   408
      HelpTag         =   ""
      Index           =   -2147483648
      InitialParent   =   ""
      Italic          =   False
      Left            =   20
      LockBottom      =   True
      LockedInPosition=   False
      LockLeft        =   True
      LockRight       =   True
      LockTop         =   True
      Panels          =   ""
      Scope           =   0
      SmallTabs       =   True
      TabDefinition   =   "About\rAP Naming and Static Channels\rRemove Overrides\rChannels\rScaling"
      TabIndex        =   11
      TabPanelIndex   =   0
      TabStop         =   True
      TextFont        =   "System"
      TextSize        =   0.0
      TextUnit        =   0
      Top             =   82
      Transparent     =   False
      Underline       =   False
      Value           =   3
      Visible         =   True
      Width           =   617
      Begin Label lblLoadCSVPrompt
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   False
         Height          =   316
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Multiline       =   False
         Scope           =   2
         Selectable      =   False
         TabIndex        =   6
         TabPanelIndex   =   2
         TabStop         =   True
         Text            =   "Please load a CSV file"
         TextAlign       =   1
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   20.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   577
      End
      Begin PushButton btnGO
         AutoDeactivate  =   True
         Bold            =   True
         ButtonStyle     =   "0"
         Cancel          =   False
         Caption         =   "GO!"
         Default         =   False
         Enabled         =   False
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   488
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Scope           =   2
         TabIndex        =   4
         TabPanelIndex   =   2
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   448
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   130
      End
      Begin PushButton btnCSVImport
         AutoDeactivate  =   True
         Bold            =   False
         ButtonStyle     =   "0"
         Cancel          =   False
         Caption         =   "Import CSV..."
         Default         =   False
         Enabled         =   True
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   False
         Scope           =   2
         TabIndex        =   1
         TabPanelIndex   =   2
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   448
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   109
      End
      Begin Label lblCSVName
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   161
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   False
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   -1
         TabPanelIndex   =   2
         TabStop         =   True
         Text            =   "No File Loaded"
         TextAlign       =   0
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   449
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   234
      End
      Begin TextArea txtAbout
         AcceptTabs      =   False
         Alignment       =   1
         AutoDeactivate  =   True
         AutomaticallyCheckSpelling=   False
         BackColor       =   &cFFFFFF00
         Bold            =   False
         Border          =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Format          =   ""
         Height          =   327
         HelpTag         =   ""
         HideSelection   =   True
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LimitText       =   0
         LineHeight      =   0.0
         LineSpacing     =   1.0
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Mask            =   ""
         Multiline       =   True
         ReadOnly        =   True
         Scope           =   0
         ScrollbarHorizontal=   False
         ScrollbarVertical=   True
         Styled          =   False
         TabIndex        =   0
         TabPanelIndex   =   1
         TabStop         =   True
         Text            =   "This tool is to be used together with a Ruckus SmartZone controller (SZ-100, SCG-200, SZ-300 or any vSZ/Virtual SmartZone).  It has been tested with SZ version 5.x upwards and uses the SZ API for all functionality.\r\rHow to use:\r“AP Naming” tab:\r- Connect to SZ [no need to select a Zone in this case]\r- Click Import CSV to import a CSV file\r- Click ‘GO!’ to begin.\r\r""Remove Overrides"" tab:\r- Connect to SZ [no need to select a Zone in this case]\r- Click Import CSV to import a CSV file\r- Select on which options you would like to remove overrides (default = all!)\r- Click 'GO!' to begin.\r\r“Channels” tab:\r- Connect to SZ\r- Select a Zone\r- In the ‘Channels’ tab, click Refresh\r- The two tables will fill up with the number of APs on each channel in that zone.\r-----\rAdditional Features:\rThe ""Export CSV"" button creates a CSV with details of APs in the selected Zone.  This CSV can be edited (in Excel, Sheets, etc) and then re-imported as is in the ""AP Naming"" screen.  \r\rYou can specify the IP, port, username and password on the command line, using the following parameters:\r-ip=, -ipaddr= or -ipaddress= to set SZ IP\r-usr= or -user= to set SZ Username\r-pwd=, or -password= to set SZ Password\r-prt=, or -port= to set SZ API Port\r-debug=yes will show a debug window\r\r-----\rFeedback welcome!\rAndrea Coppini\randrea@tacoppini.com\rhttp://dogtag.tacoppini.com\r\rThis is an unofficial, unsupported, put-together-over-several-flights tool.  I am not responsible for any outcome - good or bad!  USE AT YOUR OWN RISK\r"
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   143
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   False
         Visible         =   True
         Width           =   577
      End
      Begin PushButton btnChanLoad
         AutoDeactivate  =   True
         Bold            =   False
         ButtonStyle     =   "0"
         Cancel          =   False
         Caption         =   "Refresh"
         Default         =   False
         Enabled         =   False
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         Scope           =   2
         TabIndex        =   0
         TabPanelIndex   =   4
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   100
      End
      Begin Label lblAPListCnt
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   407
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   5
         TabPanelIndex   =   2
         TabStop         =   True
         Text            =   ""
         TextAlign       =   2
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   449
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   68
      End
      Begin Label lblAboutDogTag
         AutoDeactivate  =   True
         Bold            =   True
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Multiline       =   False
         Scope           =   0
         Selectable      =   True
         TabIndex        =   1
         TabPanelIndex   =   1
         TabStop         =   True
         Text            =   "DogTag"
         TextAlign       =   0
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   577
      End
      Begin PopupMenu popAPGrp
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   False
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   ""
         Italic          =   False
         Left            =   152
         ListIndex       =   0
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Scope           =   0
         TabIndex        =   6
         TabPanelIndex   =   4
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   False
         Underline       =   False
         Visible         =   False
         Width           =   234
      End
      Begin Listbox listChan2
         AutoDeactivate  =   True
         AutoHideScrollbars=   True
         Bold            =   False
         Border          =   True
         ColumnCount     =   2
         ColumnsResizable=   True
         ColumnWidths    =   ""
         DataField       =   ""
         DataSource      =   ""
         DefaultRowHeight=   -1
         Enabled         =   False
         EnableDrag      =   False
         EnableDragReorder=   False
         GridLinesHorizontal=   2
         GridLinesVertical=   2
         HasHeading      =   True
         HeadingIndex    =   -1
         Height          =   284
         HelpTag         =   ""
         Hierarchical    =   False
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   "Channel	APs"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         RequiresSelection=   False
         Scope           =   2
         ScrollbarHorizontal=   True
         ScrollBarVertical=   True
         SelectionType   =   1
         ShowDropIndicator=   False
         TabIndex        =   7
         TabPanelIndex   =   4
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   186
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   True
         Width           =   265
         _ScrollOffset   =   0
         _ScrollWidth    =   -1
      End
      Begin Listbox listChan5
         AutoDeactivate  =   True
         AutoHideScrollbars=   True
         Bold            =   False
         Border          =   True
         ColumnCount     =   2
         ColumnsResizable=   True
         ColumnWidths    =   ""
         DataField       =   ""
         DataSource      =   ""
         DefaultRowHeight=   -1
         Enabled         =   False
         EnableDrag      =   False
         EnableDragReorder=   False
         GridLinesHorizontal=   2
         GridLinesVertical=   2
         HasHeading      =   True
         HeadingIndex    =   -1
         Height          =   284
         HelpTag         =   ""
         Hierarchical    =   True
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   "Channel	APs"
         Italic          =   False
         Left            =   353
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         RequiresSelection=   False
         Scope           =   2
         ScrollbarHorizontal=   True
         ScrollBarVertical=   True
         SelectionType   =   1
         ShowDropIndicator=   False
         TabIndex        =   8
         TabPanelIndex   =   4
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   186
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   True
         Width           =   265
         _ScrollOffset   =   0
         _ScrollWidth    =   -1
      End
      Begin Label lblChan2G
         AutoDeactivate  =   True
         Bold            =   True
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         Multiline       =   False
         Scope           =   2
         Selectable      =   False
         TabIndex        =   9
         TabPanelIndex   =   4
         TabStop         =   True
         Text            =   "2.4Ghz"
         TextAlign       =   0
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   154
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   100
      End
      Begin Label lblChan5G
         AutoDeactivate  =   True
         Bold            =   True
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   353
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         Multiline       =   False
         Scope           =   2
         Selectable      =   False
         TabIndex        =   10
         TabPanelIndex   =   4
         TabStop         =   True
         Text            =   "5Ghz"
         TextAlign       =   0
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   154
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   100
      End
      Begin PushButton btnorideCSVImport
         AutoDeactivate  =   True
         Bold            =   False
         ButtonStyle     =   "0"
         Cancel          =   False
         Caption         =   "Import CSV..."
         Default         =   False
         Enabled         =   True
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   False
         Scope           =   2
         TabIndex        =   0
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   448
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   109
      End
      Begin PushButton btnorideGO
         AutoDeactivate  =   True
         Bold            =   True
         ButtonStyle     =   "0"
         Cancel          =   False
         Caption         =   "GO!"
         Default         =   False
         Enabled         =   False
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   490
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Scope           =   2
         TabIndex        =   2
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   448
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   130
      End
      Begin Label lblorideCSVName
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   161
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   False
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   3
         TabPanelIndex   =   3
         TabStop         =   True
         Text            =   "No File Loaded"
         TextAlign       =   0
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   449
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   234
      End
      Begin Label lblorideAPListCnt
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   413
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   4
         TabPanelIndex   =   3
         TabStop         =   True
         Text            =   ""
         TextAlign       =   2
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   450
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   68
      End
      Begin Listbox listorideOptions
         AutoDeactivate  =   True
         AutoHideScrollbars=   True
         Bold            =   False
         Border          =   True
         ColumnCount     =   2
         ColumnsResizable=   True
         ColumnWidths    =   ""
         DataField       =   ""
         DataSource      =   ""
         DefaultRowHeight=   -1
         Enabled         =   True
         EnableDrag      =   False
         EnableDragReorder=   False
         GridLinesHorizontal=   2
         GridLinesVertical=   2
         HasHeading      =   True
         HeadingIndex    =   -1
         Height          =   152
         HelpTag         =   ""
         Hierarchical    =   False
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   "Option	Remove"
         Italic          =   False
         Left            =   357
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   True
         RequiresSelection=   False
         Scope           =   2
         ScrollbarHorizontal=   True
         ScrollBarVertical=   True
         SelectionType   =   0
         ShowDropIndicator=   False
         TabIndex        =   5
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   True
         Width           =   261
         _ScrollOffset   =   0
         _ScrollWidth    =   -1
      End
      Begin GroupBox grporide2G
         AutoDeactivate  =   True
         Bold            =   False
         Caption         =   "2Ghz Radio"
         Enabled         =   True
         Height          =   90
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   356
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Scope           =   2
         TabIndex        =   6
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   347
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   125
         Begin RadioButton radoride2G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "On"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   1
            InitialParent   =   "grporide2G"
            Italic          =   False
            Left            =   427
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   0
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   373
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   46
         End
         Begin RadioButton radoride2G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Off"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   0
            InitialParent   =   "grporide2G"
            Italic          =   False
            Left            =   373
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   1
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   373
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   50
         End
         Begin RadioButton radoride2G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Skip"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   2
            InitialParent   =   "grporide2G"
            Italic          =   False
            Left            =   373
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   2
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   405
            Transparent     =   False
            Underline       =   False
            Value           =   True
            Visible         =   True
            Width           =   100
         End
      End
      Begin Label lblLoadCSVPromptORide
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   False
         Height          =   316
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Multiline       =   False
         Scope           =   2
         Selectable      =   False
         TabIndex        =   9
         TabPanelIndex   =   3
         TabStop         =   True
         Text            =   "Please load a CSV file"
         TextAlign       =   1
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   20.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   True
         Underline       =   False
         Visible         =   True
         Width           =   304
      End
      Begin GroupBox grporide5G
         AutoDeactivate  =   True
         Bold            =   False
         Caption         =   "5Ghz Radio"
         Enabled         =   True
         Height          =   89
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   493
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Scope           =   2
         TabIndex        =   7
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   347
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   125
         Begin RadioButton radoride5G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Off"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   0
            InitialParent   =   "grporide5G"
            Italic          =   False
            Left            =   505
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   0
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   374
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   48
         End
         Begin RadioButton radoride5G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "On"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   1
            InitialParent   =   "grporide5G"
            Italic          =   False
            Left            =   565
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   1
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   374
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   48
         End
         Begin RadioButton radoride5G
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Skip"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   2
            InitialParent   =   "grporide5G"
            Italic          =   False
            Left            =   505
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   2
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   406
            Transparent     =   False
            Underline       =   False
            Value           =   True
            Visible         =   True
            Width           =   108
         End
      End
      Begin GroupBox grporideIPAddr
         AutoDeactivate  =   True
         Bold            =   False
         Caption         =   "AP IP Address"
         Enabled         =   True
         Height          =   55
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   357
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   False
         Scope           =   2
         TabIndex        =   8
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   280
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   261
         Begin RadioButton radorideIPAddr
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Skip"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   1
            InitialParent   =   "grporideIPAddr"
            Italic          =   False
            Left            =   505
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   0
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   306
            Transparent     =   False
            Underline       =   False
            Value           =   True
            Visible         =   True
            Width           =   104
         End
         Begin RadioButton radorideIPAddr
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Force Dynamic"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   0
            InitialParent   =   "grporideIPAddr"
            Italic          =   False
            Left            =   374
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   2
            TabIndex        =   1
            TabPanelIndex   =   3
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   306
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   114
         End
      End
      Begin Listbox listorideAP
         AutoDeactivate  =   True
         AutoHideScrollbars=   False
         Bold            =   False
         Border          =   True
         ColumnCount     =   3
         ColumnsResizable=   True
         ColumnWidths    =   ""
         DataField       =   ""
         DataSource      =   ""
         DefaultRowHeight=   -1
         Enabled         =   True
         EnableDrag      =   False
         EnableDragReorder=   False
         GridLinesHorizontal=   2
         GridLinesVertical=   2
         HasHeading      =   True
         HeadingIndex    =   -1
         Height          =   316
         HelpTag         =   ""
         Hierarchical    =   False
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   "MAC Address	OK	Error"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         RequiresSelection=   False
         Scope           =   2
         ScrollbarHorizontal=   True
         ScrollBarVertical=   True
         SelectionType   =   1
         ShowDropIndicator=   False
         TabIndex        =   1
         TabPanelIndex   =   3
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   121
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   False
         Width           =   304
         _ScrollOffset   =   0
         _ScrollWidth    =   -1
      End
      Begin Listbox listAP
         AutoDeactivate  =   True
         AutoHideScrollbars=   True
         Bold            =   False
         Border          =   True
         ColumnCount     =   11
         ColumnsResizable=   True
         ColumnWidths    =   ""
         DataField       =   ""
         DataSource      =   ""
         DefaultRowHeight=   -1
         Enabled         =   True
         EnableDrag      =   False
         EnableDragReorder=   False
         GridLinesHorizontal=   2
         GridLinesVertical=   2
         HasHeading      =   True
         HeadingIndex    =   -1
         Height          =   316
         HelpTag         =   ""
         Hierarchical    =   False
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         InitialValue    =   "MAC Address	Name	Description	Location	Latitude	Longitude	Static Width (2.4G)	Static Channel (2.4G)	Static Width (5G)	Static Channel (5G)	Result"
         Italic          =   False
         Left            =   40
         LockBottom      =   True
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         RequiresSelection=   False
         Scope           =   2
         ScrollbarHorizontal=   True
         ScrollBarVertical=   True
         SelectionType   =   1
         ShowDropIndicator=   False
         TabIndex        =   0
         TabPanelIndex   =   2
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   121
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   False
         Width           =   577
         _ScrollOffset   =   0
         _ScrollWidth    =   -1
      End
      Begin GroupBox grpScaleConfig
         AutoDeactivate  =   True
         Bold            =   False
         Caption         =   "SmartZone Configuration"
         Enabled         =   True
         Height          =   151
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   40
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   False
         LockTop         =   True
         Scope           =   0
         TabIndex        =   0
         TabPanelIndex   =   5
         TabStop         =   True
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   120
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   254
         Begin RadioButton radScalDash
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "Essentials"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   0
            InitialParent   =   "grpScaleConfig"
            Italic          =   False
            Left            =   60
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   0
            TabIndex        =   0
            TabPanelIndex   =   5
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   156
            Transparent     =   False
            Underline       =   False
            Value           =   True
            Visible         =   True
            Width           =   100
         End
         Begin RadioButton radScalDash
            AutoDeactivate  =   True
            Bold            =   False
            Caption         =   "High Scale"
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   1
            InitialParent   =   "grpScaleConfig"
            Italic          =   False
            Left            =   172
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   0
            TabIndex        =   1
            TabPanelIndex   =   5
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   156
            Transparent     =   False
            Underline       =   False
            Value           =   False
            Visible         =   True
            Width           =   100
         End
         Begin PopupMenu popScalVer
            AutoDeactivate  =   True
            Bold            =   False
            DataField       =   ""
            DataSource      =   ""
            Enabled         =   True
            Height          =   22
            HelpTag         =   ""
            Index           =   -2147483648
            InitialParent   =   "grpScaleConfig"
            InitialValue    =   "pre 5.0\r\n5.0\r\n5.1 onwards"
            Italic          =   False
            Left            =   129
            ListIndex       =   2
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Scope           =   0
            TabIndex        =   2
            TabPanelIndex   =   5
            TabStop         =   True
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   188
            Transparent     =   False
            Underline       =   False
            Visible         =   True
            Width           =   144
         End
         Begin Label lblScalVersion
            AutoDeactivate  =   True
            Bold            =   False
            DataField       =   ""
            DataSource      =   ""
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   -2147483648
            InitialParent   =   "grpScaleConfig"
            Italic          =   False
            Left            =   60
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Multiline       =   False
            Scope           =   0
            Selectable      =   False
            TabIndex        =   3
            TabPanelIndex   =   5
            TabStop         =   True
            Text            =   "Version"
            TextAlign       =   0
            TextColor       =   &c00000000
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   189
            Transparent     =   False
            Underline       =   False
            Visible         =   True
            Width           =   66
         End
         Begin Slider sldScalClusterSize
            AutoDeactivate  =   True
            Enabled         =   True
            Height          =   24
            HelpTag         =   ""
            Index           =   -2147483648
            InitialParent   =   "grpScaleConfig"
            Left            =   129
            LineStep        =   1
            LiveScroll      =   True
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Maximum         =   4
            Minimum         =   1
            PageStep        =   1
            Scope           =   0
            TabIndex        =   4
            TabPanelIndex   =   5
            TabStop         =   True
            TickStyle       =   "1"
            Top             =   222
            Transparent     =   False
            Value           =   1
            Visible         =   True
            Width           =   143
         End
         Begin Label lblScalClusterSize
            AutoDeactivate  =   True
            Bold            =   False
            DataField       =   ""
            DataSource      =   ""
            Enabled         =   True
            Height          =   20
            HelpTag         =   ""
            Index           =   -2147483648
            InitialParent   =   "grpScaleConfig"
            Italic          =   False
            Left            =   60
            LockBottom      =   False
            LockedInPosition=   False
            LockLeft        =   True
            LockRight       =   False
            LockTop         =   True
            Multiline       =   False
            Scope           =   0
            Selectable      =   False
            TabIndex        =   5
            TabPanelIndex   =   5
            TabStop         =   True
            Text            =   "Nodes"
            TextAlign       =   0
            TextColor       =   &c00000000
            TextFont        =   "System"
            TextSize        =   0.0
            TextUnit        =   0
            Top             =   226
            Transparent     =   False
            Underline       =   False
            Visible         =   True
            Width           =   66
         End
      End
      Begin Slider sldScalAPSwitch
         AutoDeactivate  =   True
         Enabled         =   True
         Height          =   24
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Left            =   306
         LineStep        =   1
         LiveScroll      =   True
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Maximum         =   1024
         Minimum         =   0
         PageStep        =   20
         Scope           =   0
         TabIndex        =   3
         TabPanelIndex   =   5
         TabStop         =   True
         TickStyle       =   "2"
         Top             =   247
         Transparent     =   False
         Value           =   1024
         Visible         =   True
         Width           =   311
      End
      Begin TextField txtScalAP
         AcceptTabs      =   False
         Alignment       =   0
         AutoDeactivate  =   True
         AutomaticallyCheckSpelling=   False
         BackColor       =   &cFFFFFF00
         Bold            =   False
         Border          =   True
         CueText         =   ""
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Format          =   ""
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   517
         LimitText       =   0
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   True
         Mask            =   ""
         Password        =   False
         ReadOnly        =   True
         Scope           =   0
         TabIndex        =   2
         TabPanelIndex   =   5
         TabStop         =   True
         Text            =   "1024"
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   139
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   True
         Width           =   100
      End
      Begin TextField txtScalSwitch
         AcceptTabs      =   False
         Alignment       =   0
         AutoDeactivate  =   True
         AutomaticallyCheckSpelling=   False
         BackColor       =   &cFFFFFF00
         Bold            =   False
         Border          =   True
         CueText         =   ""
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Format          =   ""
         Height          =   22
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   517
         LimitText       =   0
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   True
         Mask            =   ""
         Password        =   False
         ReadOnly        =   True
         Scope           =   0
         TabIndex        =   1
         TabPanelIndex   =   5
         TabStop         =   True
         Text            =   "0"
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   173
         Transparent     =   False
         Underline       =   False
         UseFocusRing    =   True
         Visible         =   True
         Width           =   100
      End
      Begin Label lblScalAP
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   405
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   False
         LockRight       =   True
         LockTop         =   True
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   4
         TabPanelIndex   =   5
         TabStop         =   True
         Text            =   "Access Points"
         TextAlign       =   2
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   140
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   100
      End
      Begin Label lblScalSwitch
         AutoDeactivate  =   True
         Bold            =   False
         DataField       =   ""
         DataSource      =   ""
         Enabled         =   True
         Height          =   20
         HelpTag         =   ""
         Index           =   -2147483648
         InitialParent   =   "tabMain"
         Italic          =   False
         Left            =   405
         LockBottom      =   False
         LockedInPosition=   False
         LockLeft        =   True
         LockRight       =   True
         LockTop         =   True
         Multiline       =   False
         Scope           =   0
         Selectable      =   False
         TabIndex        =   5
         TabPanelIndex   =   5
         TabStop         =   True
         Text            =   "Switches"
         TextAlign       =   2
         TextColor       =   &c00000000
         TextFont        =   "System"
         TextSize        =   0.0
         TextUnit        =   0
         Top             =   174
         Transparent     =   False
         Underline       =   False
         Visible         =   True
         Width           =   100
      End
   End
End
#tag EndWindow

#tag WindowCode
	#tag Event
		Sub Open()
		  lblAboutDogTag.Text = "DogTag " + Str(App.MajorVersion) + "." + Str(App.MinorVersion) + "." + Str(App.BugVersion) + "." + Str(App.NonReleaseVersion) + " ("+ Str(App.BuildDate.AbbreviatedDate) + ")"
		  MsgBox("Note:  This version of DogTag only works with SmartZone OS v5 and above.")
		  dim CmdAll, Args(), ArgPair(1) as string
		  Args = split(system.commandLine, " ")
		  
		  For Each Arg As String In Args() 
		    If Arg.InStr("=") >0  Then
		      ArgPair = Split(Arg, "=")
		      
		      Select Case ArgPair(0)
		      Case "-ip","-ipaddr","-ipaddress"
		        txtDebug.AppendText("Found IP " + ArgPair(1) + " on command line... setting IP" + EndOfLine)
		        txtSZIP.Text = ArgPair(1)
		      Case "-usr","-user"
		        txtDebug.AppendText("Found User " + ArgPair(1)+ " on command line... setting User" + EndOfLine)
		        txtSZUName.Text = ArgPair(1)
		      Case "-pwd","-password"
		        txtDebug.AppendText("Found Password [hidden] on command line... setting Password" + EndOfLine)
		        txtSZPwd.Text = ArgPair(1)
		      Case "-prt","-port"
		        txtDebug.AppendText("Found Port " + ArgPair(1) + " on command line... setting Port" + EndOfLine)
		        txtSZPort.Text = ArgPair(1)
		      Case "-debug"
		        if (ArgPair(1) = "yes" or ArgPair(1) = "y" or ArgPair(1) = "true") Then
		          tabMain.Height = DogTag.Height - 175
		          txtDebug.Visible = True
		        End if
		      End Select
		    End If
		    
		    
		  Next
		  
		  arrorideURLList() = Array("login","syslog","wifi24","wifi50","wlanGroup24","wlanGroup50","smartMonitor","bonjourGateway","clientAdmissionControl24","clientAdmissionControl50","venueProfile","specific","channelEvaluationInterval","autoChannelSelection24","autoChannelSelection50","picture","channelEvaluationInterval","meshOptions")
		  
		  For Each URI As String in arrorideURLList ()
		    listorideOptions.AddRow(URI)
		    arrorideURL.Append(URI)
		    listorideOptions.CellType(listorideOptions.LastIndex,1) = 2
		    listorideOptions.CellCheck(listorideOptions.LastIndex,1) = True
		  Next
		End Sub
	#tag EndEvent

	#tag Event
		Sub Resized()
		  listChan2.Width = (me.Width / 2) - 60
		  listChan5.Width = (me.Width / 2) - 60
		  listChan5.Left = (me.Width / 2) + 28
		  lblChan5G.Left = (me.Width / 2) + 28
		  
		  'cha2.Width = (me.Width / 2) - 60
		  'cha5.Width = (me.Width / 2) - 60
		  'cha5.Left = (me.Width / 2) + 28
		  
		  'cha2.redraw
		  'cha5.redraw
		End Sub
	#tag EndEvent

	#tag Event
		Sub Resizing()
		  listChan2.Width = (me.Width / 2) - 60
		  listChan5.Width = (me.Width / 2) - 60
		  listChan5.Left = (me.Width / 2) + 28
		  lblChan5G.Left = (me.Width / 2) + 28
		End Sub
	#tag EndEvent


	#tag Method, Flags = &h0
		Sub ScalCalc()
		  Dim MaxSwitch, MaxAP,scalAPSWRatio, scalClusterMult, szVer As Double
		  Dim boolSwitch As Boolean
		  
		  Select Case popScalVer.ListIndex
		  Case 0  // Pre 5.0
		    txtScalSwitch.Enabled = False
		    lblScalSwitch.Enabled = False
		    boolSwitch = False
		    szVer = 3.0
		  Case 1  // 5.0
		    txtScalSwitch.Enabled = True
		    lblScalSwitch.Enabled = True
		    boolSwitch = True
		    szVer = 5.0
		  Case 2  // 5.1 onwards
		    txtScalSwitch.Enabled = True
		    lblScalSwitch.Enabled = True
		    boolSwitch = True
		    szVer = 5.1
		  End Select
		  
		  ScalClusterMult = 1
		  
		  Select Case sldScalClusterSize.Value
		  Case 1 
		    ScalClusterMult = 1
		  Case 2
		    ScalClusterMult = 1
		  Case 3
		    ScalClusterMult = 2
		  Case 4
		    ScalClusterMult = 3
		  End Select
		  
		  If radScalDash(0).Value = True Then
		    MaxAP = 1024 
		    If szVer = 5.0 Then scalAPSWRatio = 20
		    If szVer = 5.1 Then scalAPSWRatio = 5
		  End If
		  
		  If radScalDash(1).Value = True Then
		    MaxAP = 10000
		    If szVer = 5.0 Then scalAPSWRatio = 20
		    If szVer = 5.1 Then scalAPSWRatio = 8
		  End If
		  
		  MaxAP = MaxAP * ScalClusterMult
		  sldScalAPSwitch.Maximum = MaxAP
		  
		  If boolSwitch Then
		    sldScalAPSwitch.Enabled = True
		    MaxSwitch = MaxAP / scalAPSWRatio
		    txtScalAP.text = Str(sldScalAPSwitch.Value)
		    txtScalSwitch.Text = Str(Floor(MaxSwitch - (sldScalAPSwitch.Value / scalAPSWRatio)))
		  Else
		    sldScalAPSwitch.Value = MaxAP
		    sldScalAPSwitch.Enabled = False
		    txtScalAP.text = Str(sldScalAPSwitch.Value)
		    txtScalSwitch.Text = "N/A"
		  End If
		End Sub
	#tag EndMethod

	#tag Method, Flags = &h0
		Function SZGet(GetURL As String,GetBody As String) As String
		  timeKeepalive.Reset
		  SZ.Secure = True
		  SZ.Port = Val(txtSZPort.Text)
		  SZ.Yield = True
		  SZ.ConnectionType = SSLSocket.TLSv12
		  SZ.SetRequestHeader("Content-Type", "application/json;charset=UTF-8")
		  SZURL = "https://" + txtSZIP.Text + ":" + txtSZPort.Text + "wsg/api/public"
		  
		  If GetBody <> "" Then
		    SZ.SetRequestContent(GetBody, "application/json;charset=UTF-8") 
		    txtDebug.AppendText("Adding Body: " + GetBody + EndOfLine)
		  End If
		  
		  If SZCookie <> "" Then
		    SZ.SetRequestHeader("Cookie", SZCookie)
		    Return SZ.Get(SZURL+GetURL, 5)
		  Else
		    MsgBox("Error during GET request - Session not initialised")
		  End If
		  timeKeepalive.Reset
		  
		End Function
	#tag EndMethod

	#tag Method, Flags = &h0
		Sub SZLogoff()
		  Dim data As String
		  
		  SZ.SetRequestHeader("Cookie", SZCookie)
		  SZ.SetRequestContent("", "application/json;charset=UTF-8") 
		  data = SZ.SendRequest("DELETE", SZURL+"/v7_0/session",5)
		  If SZ.ErrorCode <> -1 Then 
		    SZCookie = ""
		    btnConnect.Caption = "Connect"
		    ovlConnect.FillColor = Red
		    popZone.DeleteAllRows
		    popZone.Enabled = False
		    btnGo.Enabled=False
		    btnorideGo.Enabled=False
		    btnChanLoad.Enabled=False
		    btnCSVExport.Enabled=False
		  End if
		End Sub
	#tag EndMethod

	#tag Method, Flags = &h0
		Function SZPost() As String
		  
		End Function
	#tag EndMethod


#tag EndWindowCode

#tag Events txtDebug
	#tag Event
		Sub TextChange()
		  txtDebug.ScrollPosition = 9999
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnConnect
	#tag Event
		Sub Action()
		  Dim data, dataError, SZContent As String
		  Dim SZHeaders As InternetHeaders
		  Dim SZZoneID(), SZZoneName() As String
		  
		  App.DoEvents
		  
		  SZ.Secure = True
		  SZ.Port = Val(txtSZPort.Text)
		  SZ.Yield = True
		  SZ.ConnectionType = SSLSocket.TLSv12
		  SZ.SetRequestHeader("Content-Type", "application/json;charset=UTF-8")
		  SZURL = "https://" + txtSZIP.Text + ":" + txtSZPort.Text + "wsg/api/public"
		  
		  if SZCookie = "" Then 
		    ' Logon
		    If txtSZIP.Text <> "" and txtSZPort.Text <> "" and txtSZUName.Text <> "" and txtSZPwd.Text <> "" Then
		      SZContent = "{""username"": """+txtSZUName.Text+""",""password"": """+txtSZPwd.Text+""",""timeZoneUtcOffset"": ""+02:00""}"
		      SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		      data = SZ.Post(SZURL+"/v7_0/session", 5)  
		      Dim dataResponse As New JSONItem(data)
		      dataError = dataResponse.Lookup("controllerVersion","Error connecting to SZ")
		      SZVer = dataResponse.Lookup("controllerVersion","0")
		      txtDebug.AppendText("Detected version: " + SZVer + EndOfLine)
		      If dataError <> "Error connecting to SZ" Then
		        If SZ.ErrorCode <> -1 Then 
		          SZHeaders = SZ.PageHeaders
		          SZCookie = SZHeaders.Value("Set-cookie")
		          
		          btnConnect.Caption = "Disconnect"
		          ovlConnect.FillColor = Green
		          
		          txtDebug.AppendText ("Connected... Retrieving Zone List" + EndOfLine)
		          popZone.DeleteAllRows
		          
		          SZZones = SZGet("/v7_0/rkszones?listSize=1000","")
		          
		          Dim SZZoneAll As New JSONItem(SZZones)
		          txtDebug.AppendText ("JSON: " + SZZoneAll.ToString + EndOfLine)
		          
		          Dim SZZoneList As JSONItem = SZZoneAll.Value("list")
		          Redim SZZoneArr(SZZoneList.Count-1,3) ' Resize the global array
		          popZone.DeleteAllRows
		          
		          For i As Integer = 0 to SZZoneList.Count -1 
		            Dim SZZoneItem As JSONItem = SZZoneList.Child(i)
		            SZZoneID.Append(SZZoneItem.Value("id"))
		            SZZoneName.Append(SZZoneItem.Value("name"))
		            
		            txtDebug.AppendText (SZZoneID(i) + " = " + SZZoneName(i) + EndOfLine)
		            
		          Next
		          
		          SZZoneName.SortWith(SZZoneID)
		          
		          For i As Integer = 0 to SZZoneName.Ubound 
		            popZone.Addrow( SZZoneName(i) ) 
		            popZone.RowTag( popZone.ListCount - 1 ) = SZZoneID(i) 
		          Next i 
		          
		          popZone.PopSelByText("Default Zone")
		          popZone.PopSelByText("Staging Zone")
		          popZone.Enabled = True
		          
		          btnGo.Enabled = True
		          btnorideGo.Enabled = True
		          btnChanLoad.Enabled = True
		          btnCSVExport.Enabled = True
		        Else
		          MsgBox ("Error: " + Str(SZ.ErrorCode))
		        End If
		      Else
		        MsgBox ("Error connecting to SZ")
		      End If
		    Else
		      MsgBox ("Please enter SZ details")
		    End If
		  Else
		    SZLogoff
		  End If
		  
		  If SZ.ErrorCode = -1 Then
		    txtDebug.AppendText("Error!"+EndOfLine+Str(SZ.LastErrorCode)+EndOfLine)
		  Else
		    txtDebug.AppendText (data + EndOfLine)
		  End If
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events popZone
	#tag Event
		Sub Change()
		  SZZoneSel = popZone.RowTag(popZone.ListIndex)
		  txtDebug.AppendText (SZZoneSel + EndOfLine)
		  
		  
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events chkHidePw
	#tag Event
		Sub Action()
		  Select Case chkHidePw.State
		  Case CheckBox.CheckedStates.Unchecked
		    txtSZPwd.Password = False
		    
		  Case CheckBox.CheckedStates.Checked
		    txtSZPwd.Password = True
		    
		  Case CheckBox.CheckedStates.Indeterminate
		    // The CheckBox state is indeterminate
		    
		  End Select
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnCSVExport
	#tag Event
		Sub Action()
		  Dim csvOut As TextOutputStream
		  Dim dlg as New SaveAsDialog
		  Dim f as FolderItem
		  Dim APArr(-1,-1),data As String
		  
		  
		  dlg.InitialDirectory=SpecialFolder.Desktop
		  dlg.promptText="Save Text As"
		  dlg.SuggestedFileName="DogTag-" + txtSZIP.Text + "-" + popZone.Text + ".csv"
		  dlg.Title="Text Output"
		  
		  f=dlg.ShowModal()
		  
		  If f <> Nil then 'Proceed unless user cancelled the dialog
		    txtDebug.AppendText("Loading APs from " + SZZoneSel + EndOfLine)
		    data = SZGet("/v7_0/aps?listSize=10000&zoneId="+SZZoneSel,"")
		    
		    Dim SZAPAll As New JSONItem(data)
		    
		    Me.Enabled = False
		    If SZAPAll.ToString <> "{}" Then
		      If SZAPAll.Value("totalCount") > 0 Then
		        csvOut = TextOutputStream.Create(f)
		        Dim SZAPList As JSONItem = SZAPAll.Value("list")
		        Redim APArr(SZAPList.Count-1,15) ' Resize the global array
		        For i As Integer = 0 to SZAPList.Count -1 
		          Dim SZAPItem As JSONItem = SZAPList.Child(i)
		          Dim APMac As String
		          APArr(i,0) = SZAPItem.Value("mac")
		          APArr(i,1) = SZAPItem.Value("name")
		          txtDebug.AppendText(APArr(i,0) + EndOfLine)
		        Next
		        txtDebug.AppendText("Retrieved " + Str(SZAPList.Count) + " APs... Loading AP Data" + EndOfLine)
		        data = ""
		        
		        
		        progBar.Value = 0
		        progBar.Maximum = UBound(APArr)+1
		        progBar.Visible = True
		        For j As Integer = 0 to UBound(APArr)
		          data = SZGet("/v7_0/aps/"+APArr(j,0)+"/operational/summary","")
		          Dim APDataAll As New JSONItem(data)
		          APArr(j,2) = APDataAll.Value("description")
		          APArr(j,3) = APDataAll.Value("location")
		          APArr(j,4) = APDataAll.Value("latitude")
		          APArr(j,5) = APDataAll.Value("longitude")
		          APArr(j,6) = ""  // Channel width for 2.4Ghz (Cannot retrieve via API)
		          APArr(j,7) = APDataAll.Value("wifi24Channel")
		          APArr(j,8) = "" // Channel width for 5Ghz (Cannot retrieve via API)
		          APArr(j,9) = APDataAll.Value("wifi50Channel")
		          APArr(j,10) = APDataAll.Value("serial")
		          APArr(j,11) = APDataAll.Value("connectionState")
		          APArr(j,12) = APDataAll.Value("model")
		          APArr(j,13) = APDataAll.Value("ipType")
		          APArr(j,14) = APDataAll.Value("ip")
		          txtDebug.AppendText(APArr(j,0) + " - " + APArr(j,1) + " / " + ApArr(j,2) + " (" + ApArr(j,3) + ")" + EndOfLine)
		          
		          progBar.Value = j+1
		        Next
		        
		        if csvOut <> nil then
		          txtDebug.AppendText("Writing.." + EndOfLine)
		          csvOut.Write """MAC"",""Name"",""Description"",""Location"",""GPSLatitude"",""GPSLongitude"",""24Width"",""24Channel"",""50Width"",""50Channel"",""Serial"",""Status"",""Model"",""IPType"",""IPAddress""" + EndOfLine
		          For k As Integer = 0 to UBound(APArr)
		            For l As Integer = 0 to UBound(APArr,2)
		              csvOut.Write """" + APArr(k,l) + ""","
		            Next
		            csvOut.Write EndOfLine
		          Next
		          
		          txtDebug.AppendText("Done!" + EndOfLine)
		          progBar.Visible = False
		        Else
		          MsgBox ("Error! This Zone has no APs!")
		        End If
		        csvOut.Close
		      Else
		        MsgBox ("Error! (SZ returned empty values)")
		      End If
		      Me.Enabled = True
		      
		      
		      
		    end if
		  Else
		    //user canceled
		  End if
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events timeKeepalive
	#tag Event
		Sub Action()
		  If SZCookie <> "" Then
		    txtDebug.AppendText ("Keep alive...")
		    SZSessionAll = SZGet("/v7_0/session","")
		    If SZSessionAll = "{ ""success"" : false }" Then
		      txtDebug.AppendText ("Session ended by server, disconnecting..." + EndOfLine)
		      SZLogoff
		    Else
		      txtDebug.AppendText (SZSessionAll + EndOfLine)
		    End If
		  Else
		    txtDebug.AppendText ("No Keep alive" + EndOfLine)
		  End If
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnGO
	#tag Event
		Sub Action()
		  Dim data,SZContent,apMAC,apName,apDesc,apLoc,apLat,apLon,ap24Width,ap24Chan,ap50Width,ap50Chan As String
		  SZ.SetRequestHeader("Cookie", SZCookie)
		  Me.Enabled = False
		  
		  If listAP.ListCount > 0 Then
		    
		    progBar.Value = 0
		    progBar.Maximum = listAP.listcount
		    progBar.Visible = True
		    for i as integer=0 to listAP.listcount-1 //this loops through all of the rows
		      apMAC = listAP.cell(i,0)
		      apName =listAP.cell(i,1)
		      apDesc = listAP.cell(i,2)
		      apLoc = listAP.cell(i,3)
		      apLat = listAP.cell(i,4)
		      apLon = listAP.cell(i,5) 
		      ap24Width = listAP.cell(i,6)
		      ap24Chan = listAP.cell(i,7)
		      ap50Width = listAP.cell(i,8)
		      ap50Chan = listAP.cell(i,9)
		      
		      SZContent = "{""name"": """+apName+""""
		      If Trim(apDesc)<>"" Then SZContent = SZContent + ",""description"":"""+apDesc+""""
		      If Trim(apLoc)<>"" Then SZContent = SZContent + ",""location"":"""+apLoc+""""
		      If Len(Trim(apLon)+Trim(apLat))>2 Then SZContent = SZContent + ",""latitude"":"""+apLat+""",""longitude"":"""+apLon+""""
		      If Len(Trim(ap24Width)+Trim(ap24Chan))>2 Then SZContent = SZContent + ",""wifi24"": {""channelWidth"": " + ap24Width + ", ""channel"": " + ap24Chan + "}"
		      If Len(Trim(ap50Width)+Trim(ap50Chan))>3 Then SZContent = SZContent + ",""wifi50"": {""channelWidth"": " + ap50Width + ", ""channel"": " + ap50Chan + "}"
		      SZContent = SZContent + "}"
		      
		      txtDebug.AppendText("Sending: "+SZContent)
		      SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		      data = SZ.SendRequest("PATCH", SZURL+"/v7_0/aps/"+apMAC,5)
		      
		      If data = "" Then
		        txtdebug.AppendText("...OK!" + EndOfLine)
		        listAP.cell(i,10) = "OK!"
		      Else
		        txtdebug.AppendText(" - ERROR" + EndOfLine+ data + EndOfLine)
		        listAP.cell(i,10) = "ERROR! - " + data
		      End If
		      
		      SZContent = ""
		      progBar.Value = i+1
		    next
		    progBar.Visible = False
		    txtdebug.AppendText("All done!" + EndOfLine + "--" + EndOfLine)
		  Else
		    txtdebug.AppendText("No APs to import!  Please select a CSV file." + EndOfLine + "--" + EndOfLine)
		  End If
		  Me.Enabled = True
		  
		  
		  '
		  
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnCSVImport
	#tag Event
		Sub Action()
		  Dim f As FolderItem = GetOpenFolderItem("text") // as defined in File Type Sets Editor
		  Dim tFile, arrFile(), arrLine() As String
		  Dim cntEntry As Integer = 0
		  
		  If f <> Nil Then
		    If f.Exists Then
		      // Be aware that TextInputStream.Open coud raise an exception
		      Dim t As TextInputStream
		      Try
		        t = TextInputStream.Open(f)
		        t.Encoding =  Encodings.UTF8
		        tFile = t.ReadAll
		        tFile = ReplaceLineEndings(tFile, EndOfLine) ' Handle different types of EOL encodings across platforms
		        arrFile = tFile.Split(EndOfLine)
		        listAP.DeleteAllRows
		        lblCSVName.Text = f.Name
		        Dim rg As New RegEx
		        Dim myMatch As RegExMatch
		        rg.SearchPattern = """?([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})""?$" 'RegEx for MAC Address
		        
		        For i As Integer = 0 To arrFile.UBound
		          ' Add to a table or list
		          If Trim(arrFile(i)) <> "" Then
		            arrLine = arrFile(i).Split(",")
		            
		            myMatch = rg.Search(arrLine(0)) 'check line for MAC Address, otherwise skip
		            If myMatch <> Nil Then
		              listAP.AddRow(arrLine)
		              listAP.Cell(listAP.LastIndex,listAP.ColumnCount-1) = "" ' Remove any imported text from the Result column
		              cntEntry = cntEntry + 1
		              For j As Integer = 0 to listAP.ColumnCount    ' Remove leading and trailing quotes
		                If Left(listAP.Cell(listAP.LastIndex, j),1) = """" Then listAP.Cell(listAP.LastIndex,j) = Right(listAP.Cell(listAP.LastIndex,j),Len(listAP.Cell(listAP.LastIndex,j))-1)
		                If Right(listAP.Cell(listAP.LastIndex, j),1) = """" Then listAP.Cell(listAP.LastIndex,j) = Left(listAP.Cell(listAP.LastIndex,j),Len(listAP.Cell(listAP.LastIndex,j))-1)
		              Next
		              
		              if Trim(listAp.Cell(listAP.LastIndex,6)) = "" and Val(Trim(listAp.Cell(listAP.LastIndex,7))) > 0 and Val(Trim(listAp.Cell(listAP.LastIndex,7))) < 15 Then 
		                listAp.Cell(listAP.LastIndex,6) = "20"
		              End If
		              if Trim(listAp.Cell(listAP.LastIndex,8)) = "" and Val(Trim(listAp.Cell(listAP.LastIndex,9))) > 30 and Val(Trim(listAp.Cell(listAP.LastIndex,9))) < 200 Then
		                listAp.Cell(listAP.LastIndex,8) = "20"
		              End If 
		              
		            End If
		          End If
		        Next
		        lblCSVName.Text = f.Name
		        lblAPListCnt.Text = Str(cntEntry)+" APs"
		        listAP.Visible=True
		        btnGo.SetFocus
		      Catch e As IOException
		        t.Close
		        MsgBox("Error accessing file.")
		      End Try
		    End If
		  End If
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnChanLoad
	#tag Event
		Sub Action()
		  Dim APArr(-1,-1),data As String
		  
		  txtDebug.AppendText("Loading APs from " + SZZoneSel + EndOfLine)
		  data = SZGet("/v7_0/aps?listSize=10000&zoneId="+SZZoneSel,"")
		  
		  Dim SZAPAll As New JSONItem(data)
		  
		  Me.Enabled = False
		  If SZAPAll.ToString <> "{}" Then
		    If SZAPAll.Value("totalCount") > 0 Then
		      listChan2.AddRow("--Started: " + NowDate + " " + NowTime + "--")
		      listChan5.AddRow("--Started:  " + NowDate + " " + NowTime + "--")
		      
		      Dim SZAPList As JSONItem = SZAPAll.Value("list")
		      Redim APArr(SZAPList.Count-1,4) ' Resize the global array
		      For i As Integer = 0 to SZAPList.Count -1 
		        Dim SZAPItem As JSONItem = SZAPList.Child(i)
		        Dim APMac As String
		        APArr(i,0) = SZAPItem.Value("mac")
		        txtDebug.AppendText(APArr(i,0) + EndOfLine)
		      Next
		      txtDebug.AppendText("Retrieved " + Str(SZAPList.Count) + " APs... Loading AP Data" + EndOfLine)
		      data = ""
		      
		      progBar.Value = 0
		      progBar.Maximum = UBound(APArr)+1
		      progBar.Visible = True
		      For j As Integer = 0 to UBound(APArr)
		        data = SZGet("/v7_0/aps/"+APArr(j,0)+"/operational/summary","")
		        Dim APDataAll As New JSONItem(data)
		        APArr(j,1) = APDataAll.Value("wifi24Channel")
		        APArr(j,2) = APDataAll.Value("wifi50Channel")
		        APArr(j,3) = APDataAll.Value("connectionState")
		        APArr(j,4) = APDataAll.Value("apGroupId")
		        txtDebug.AppendText(APArr(j,0) + " - " + APArr(j,1) + " / " + ApArr(j,2) + " (" + ApArr(j,3) + ")" + EndOfLine)
		        progBar.Value = j+1
		      Next
		      txtDebug.AppendText("Done!" + EndOfLine + "Counting Totals..." + EndOfLine)
		      
		      dim APChanDict As new Dictionary
		      dim APGrpDict As new Dictionary
		      dim ch2,ch5 As Integer
		      dim apgrp As String
		      
		      for k As integer = 0 to UBound(APArr)
		        if APArr(k,3) = "Connect" Then   ' only connected APs
		          ch2 = Val(APArr(k,1))
		          ch5 = Val(APArr(k,2))
		          apgrp = APArr(k,4)
		          APChanDict.Value(ch2) = APChanDict.Lookup(ch2, 0) + 1
		          APChanDict.Value(ch5) = APChanDict.Lookup(ch5, 0) + 1
		          APGrpDict.Value(apgrp) = APGrpDict.Lookup(apgrp, 0) + 1
		        End If    ' only connected APs
		      next
		      
		      ReDim APGroupArr(APGrpDict.Count-1,2)
		      
		      popAPGrp.AddRow ("-All-")
		      For t As Integer = 0 to APGrpDict.Count -1
		        APGroupArr(t,0) = APGrpDict.Key(t).StringValue
		        data = SZGet("/v7_0/rkszones/"+AppVars.SZZoneSel+"/apgroups/"+APGroupArr(t,0) ,"")
		        Dim APDataAll As New JSONItem(data)
		        APGroupArr(t,1) = APDataAll.Value("name")
		        popAPGrp.AddRow(APGroupArr(t,1))
		        t = t+1
		      Next
		      
		      
		      Dim ChanNumArr() As Integer
		      FOR Each key As Variant In APChanDict.Keys
		        ChanNumArr.Append(key)
		      NEXT
		      ChanNumArr.Sort
		      
		      
		      For i As integer = 0 to Ubound(ChanNumArr)
		        txtDebug.AppendText (Str(ChanNumArr(i)) + " " + APChanDict.Lookup(ChanNumArr(i),0) + EndOfLine)
		        If ChanNumArr(i) < 16 Then
		          listChan2.AddRow(Str(ChanNumArr(i)),APChanDict.Lookup(ChanNumArr(i),0))
		        End If
		        If ChanNumArr(i) > 15 Then
		          listChan5.AddRow(Str(ChanNumArr(i)),APChanDict.Lookup(ChanNumArr(i),0))
		        End If
		      Next
		      
		      listChan2.Enabled = True    
		      listChan5.Enabled = True
		      
		      listChan2.AddRow("--Ended: " + NowDate + " " + NowTime + "--")
		      listChan5.AddRow("--Ended: " + NowDate + " " + NowTime + "--")
		      
		      progBar.Visible = False
		    Else
		      MsgBox ("Error! This Zone has no APs!")
		    End If
		  Else
		    MsgBox ("Error! (SZ returned empty values)")
		  End If
		  Me.Enabled = True
		  
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events popAPGrp
	#tag Event
		Sub Change()
		  APGroupSel = APGroupArr(popAPGrp.ListIndex,1)
		  txtDebug.AppendText (APGroupSel + EndOfLine)
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events listChan2
	#tag Event
		Function KeyDown(Key As String) As Boolean
		  dim a As integer 
		  a=asc(key) 
		  
		  if a=8 or a=127 then  
		    for i as integer=listAP.listcount-1 downto 0
		      if listAP.selected(i) then
		        listAP.removeRow(i)
		      end if
		    next
		    lblAPListCnt.Text = Str(listAP.ListCount)+" APs"
		  end if 
		  
		  
		  
		End Function
	#tag EndEvent
	#tag Event
		Function CompareRows(row1 as Integer, row2 as Integer, column as Integer, ByRef result as Integer) As Boolean
		  Select Case column
		  Case 0,1
		    dim a1, a2 As Double
		    a1 = CDbl(me.Cell(row1,column))
		    a2 = CDbl(me.Cell(row2,column))
		    
		    if a1<a2 then
		      result = -1
		    elseif a1>a2 then
		      result = 1
		    else
		      result = 0
		    end if
		    return true //this means you actually did a comparison
		    
		  Else
		    return False //use the default sort
		  End Select
		End Function
	#tag EndEvent
	#tag Event
		Function HeaderPressed(column as Integer) As Boolean
		  Me.ColumnSortDirection(column) = 0
		  me.SortedColumn = -1 //----- remove sort arrows step 1
		  me.HasHeading = False //----- remove sort arrows step 2
		  me.HasHeading = True //----- remove sort arrows step 3
		End Function
	#tag EndEvent
#tag EndEvents
#tag Events listChan5
	#tag Event
		Function KeyDown(Key As String) As Boolean
		  dim a As integer 
		  a=asc(key) 
		  
		  if a=8 or a=127 then  
		    for i as integer=listAP.listcount-1 downto 0
		      if listAP.selected(i) then
		        listAP.removeRow(i)
		      end if
		    next
		    lblAPListCnt.Text = Str(listAP.ListCount)+" APs"
		  end if 
		  
		  
		  
		End Function
	#tag EndEvent
	#tag Event
		Function CompareRows(row1 as Integer, row2 as Integer, column as Integer, ByRef result as Integer) As Boolean
		  Select Case column
		  Case 0,1
		    dim a1, a2 As Double
		    a1 = CDbl(me.Cell(row1,column))
		    a2 = CDbl(me.Cell(row2,column))
		    
		    if a1<a2 then
		      result = -1
		    elseif a1>a2 then
		      result = 1
		    else
		      result = 0
		    end if
		    return true //this means you actually did a comparison
		    
		  Else
		    return False //use the default sort
		  End Select
		End Function
	#tag EndEvent
	#tag Event
		Function HeaderPressed(column as Integer) As Boolean
		  Me.ColumnSortDirection(column) = 0
		  me.SortedColumn = -1 //----- remove sort arrows step 1
		  me.HasHeading = False //----- remove sort arrows step 2
		  me.HasHeading = True //----- remove sort arrows step 3
		End Function
	#tag EndEvent
#tag EndEvents
#tag Events btnorideCSVImport
	#tag Event
		Sub Action()
		  Dim f As FolderItem = GetOpenFolderItem("text") // as defined in File Type Sets Editor
		  Dim tFile, arrFile(), arrLine() As String
		  Dim cntEntry As Integer = 0
		  
		  If f <> Nil Then
		    If f.Exists Then
		      // Be aware that TextInputStream.Open coud raise an exception
		      Dim t As TextInputStream
		      Try
		        t = TextInputStream.Open(f)
		        t.Encoding =  Encodings.UTF8
		        tFile = t.ReadAll
		        tFile = ReplaceLineEndings(tFile, EndOfLine) ' Handle different types of EOL encodings across platforms
		        arrFile = tFile.Split(EndOfLine)
		        listorideAP.DeleteAllRows
		        lblorideCSVName.Text = f.Name
		        Dim rg As New RegEx
		        Dim myMatch As RegExMatch
		        rg.SearchPattern = """?([0-9A-Fa-f]{2}[:]){5}([0-9A-Fa-f]{2})""?$" 'RegEx for MAC Address
		        
		        For i As Integer = 0 To arrFile.UBound
		          ' Add to a table or list
		          If Trim(arrFile(i)) <> "" Then
		            arrLine = arrFile(i).Split(",")
		            
		            myMatch = rg.Search(arrLine(0)) 'check line for MAC Address, otherwise skip
		            If myMatch <> Nil Then
		              listorideAP.AddRow(arrLine(0))
		              listorideAP.Cell(listorideAP.LastIndex,1) = "" ' Remove any imported text from the Result column
		              cntEntry = cntEntry + 1
		              
		              For j As Integer = 0 to listorideAP.ColumnCount    ' Remove leading and trailing quotes
		                If Left(listorideAP.Cell(listorideAP.LastIndex, j),1) = """" Then listorideAP.Cell(listorideAP.LastIndex,j) = Right(listorideAP.Cell(listorideAP.LastIndex,j),Len(listorideAP.Cell(listorideAP.LastIndex,j))-1)
		                If Right(listorideAP.Cell(listorideAP.LastIndex, j),1) = """" Then listorideAP.Cell(listorideAP.LastIndex,j) = Left(listorideAP.Cell(listorideAP.LastIndex,j),Len(listorideAP.Cell(listorideAP.LastIndex,j))-1)
		              Next
		              
		            End If
		          End If
		        Next
		        lblorideCSVName.Text = f.Name
		        lblorideAPListCnt.Text = Str(cntEntry)+" APs"
		        listorideAP.Visible=True
		        btnorideGo.SetFocus
		      Catch e As IOException
		        t.Close
		        MsgBox("Error accessing file.")
		      End Try
		    End If
		  End If
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events btnorideGO
	#tag Event
		Sub Action()
		  Dim data,SZContent,apMAC As String
		  'Dim arrorideURL() As String = Array("login","syslog","wifi24","wifi50","wlanGroup24","wlanGroup50","smartMonitor","bonjourGateway","clientAdmissionControl24","clientAdmissionControl50","venueProfile","specific","channelEvaluationInterval","autoChannelSelection24","autoChannelSelection50","meshOptions")
		  
		  SZ.SetRequestHeader("Cookie", SZCookie)
		  Me.Enabled = False
		  
		  If listorideAP.ListCount > 0 Then
		    
		    progBar.Value = 0
		    progBar.Maximum = listorideAP.listcount
		    progBar.Visible = True
		    for i as integer=0 to listorideAP.listcount-1 //this loops through all of the rows
		      apMAC = listorideAP.cell(i,0)
		      
		      If arrorideURL.Ubound > -1 Then
		        For Each URI As String in arrorideURL
		          txtDebug.AppendText("Sending: "+URI)
		          SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		          data = SZ.SendRequest("DELETE", SZURL+"/v7_0/aps/"+apMAC+"/"+URI,5)
		          
		          If data = "" Then
		            txtdebug.AppendText(" (OK!) .. ")
		            listorideAP.cell(i,1) = listorideAP.cell(i,1) + URI + ", "
		          Else
		            txtdebug.AppendText(" (ERROR: " + data + ") .. ")
		            listorideAP.cell(i,2) = listorideAP.cell(i,2) + URI + " (" + data + "), "
		          End If
		        Next
		      End If
		      
		      If radoride2G(0).Value = True Then SZContent = "{""wlanService24Enabled"": false}"
		      If radoride2G(1).Value = True Then SZContent = "{""wlanService24Enabled"": true}"
		      If radoride2G(2).Value = True Then SZContent = ""
		      
		      If SZContent <> "" Then
		        txtDebug.AppendText("Sending: "+SZContent)
		        SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		        data = SZ.SendRequest("PATCH", SZURL+"/v7_0/aps/"+apMAC,5)
		        
		        If data = "" Then
		          txtdebug.AppendText("...OK!" + EndOfLine)
		          listorideAP.cell(i,1) = listorideAP.cell(i,1) + "2G, "
		        Else
		          txtdebug.AppendText(" - 2G ERROR" + EndOfLine+ data + EndOfLine)
		          listorideAP.cell(i,2) = listorideAP.cell(i,2) + "2G (" + data + "), "
		        End If
		      End If
		      
		      If radoride5G(0).Value = True Then SZContent = "{""wlanService50Enabled"": false}"
		      If radoride5G(1).Value = True Then SZContent = "{""wlanService50Enabled"": true}"
		      If radoride5G(2).Value = True Then SZContent = ""
		      
		      If SZContent <> "" Then
		        txtDebug.AppendText("Sending: "+SZContent)
		        SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		        data = SZ.SendRequest("PATCH", SZURL+"/v7_0/aps/"+apMAC,5)
		        
		        If data = "" Then
		          txtdebug.AppendText("...OK!" + EndOfLine)
		          listorideAP.cell(i,1) = listorideAP.cell(i,1) + "5G, "
		        Else
		          txtdebug.AppendText(" - 5G ERROR" + EndOfLine+ data + EndOfLine)
		          listorideAP.cell(i,2) = listorideAP.cell(i,2) + "5G (" + data + "), "
		        End If
		      End If
		      
		      
		      If radorideIPAddr(0).Value = True Then SZContent = "{""ipType"": ""Dynamic""}"
		      If radorideIPAddr(1).Value = True Then SZContent = ""
		      
		      If SZContent <> "" Then
		        txtDebug.AppendText("Sending: "+SZContent)
		        SZ.SetRequestContent(SZContent, "application/json;charset=UTF-8") 
		        data = SZ.SendRequest("PATCH", SZURL+"/v7_0/aps/"+apMAC+"/network",5)
		        
		        If data = "" Then
		          txtdebug.AppendText("...OK!" + EndOfLine)
		          listorideAP.cell(i,1) = listorideAP.cell(i,1) + "IP, "
		        Else
		          txtdebug.AppendText(" - 5G ERROR" + EndOfLine+ data + EndOfLine)
		          listorideAP.cell(i,2) = listorideAP.cell(i,2) + "IP (" + data + "), "
		        End If
		      End If
		      
		      txtDebug.AppendText(EndOfLine)
		      
		      progBar.Value = i+1
		    next
		    progBar.Visible = False
		    txtdebug.AppendText("All done!" + EndOfLine + "--" + EndOfLine)
		  Else
		    txtdebug.AppendText("No APs to import!  Please select a CSV file." + EndOfLine + "--" + EndOfLine)
		  End If
		  Me.Enabled = True
		  
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events listorideOptions
	#tag Event
		Function KeyDown(Key As String) As Boolean
		  'dim a As integer 
		  'a=asc(key) 
		  
		  'if a=8 or a=127 then  
		  'for i as integer=listorideOptions.listcount-1 downto 0
		  'if listorideOptions.selected(i) then
		  'listorideOptions.removeRow(i)
		  'end if
		  'next
		  'end if 
		  
		  
		  
		  
		End Function
	#tag EndEvent
	#tag Event
		Sub CellAction(row As Integer, column As Integer)
		  If column = 1 Then // Is this the CheckBox column?
		    If Me.CellCheck(row, column) Then
		      arrorideURL.Append(Me.Cell(row,0))
		    Else
		      arrorideURL.Remove(arrorideURL.IndexOf(Me.Cell(row,0)))
		    End If
		    
		    // Debug section
		    txtDebug.AppendText ("List Array: ("+arrorideURL.Ubound.ToText+")")
		    For i As Integer = 0 To arrorideURL.Ubound
		      txtDebug.AppendText(arrorideURL(i) + ",")
		    Next
		    txtDebug.AppendText(EndOfLine)
		    // End Debug section
		    
		  End If
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events listorideAP
	#tag Event
		Function KeyDown(Key As String) As Boolean
		  dim a As integer 
		  a=asc(key) 
		  
		  if a=8 or a=127 then  
		    for i as integer=listorideAP.listcount-1 downto 0
		      if listorideAP.selected(i) then
		        listorideAP.removeRow(i)
		      end if
		    next
		    lblorideAPListCnt.Text = Str(listorideAP.ListCount)+" APs"
		  end if 
		  
		  
		  
		End Function
	#tag EndEvent
#tag EndEvents
#tag Events listAP
	#tag Event
		Function KeyDown(Key As String) As Boolean
		  dim a As integer 
		  a=asc(key) 
		  
		  if a=8 or a=127 then  
		    for i as integer=listAP.listcount-1 downto 0
		      if listAP.selected(i) then
		        listAP.removeRow(i)
		      end if
		    next
		    lblAPListCnt.Text = Str(listAP.ListCount)+" APs"
		  end if 
		  
		  
		  
		End Function
	#tag EndEvent
	#tag Event
		Sub Open()
		  Dim colWidth As String
		  
		  colWidth = Str((listAP.Width / listAP.ColumnCount) + 20)
		  
		  
		  For j As integer = 0 to listAP.ColumnCount 
		    colWidth = colWidth + "," + colWidth
		  next 
		  
		  listAP.ColumnWidths = colWidth
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events radScalDash
	#tag Event
		Sub Action(index as Integer)
		  ScalCalc
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events popScalVer
	#tag Event
		Sub Change()
		  ScalCalc
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events sldScalClusterSize
	#tag Event
		Sub ValueChanged()
		  ScalCalc
		End Sub
	#tag EndEvent
#tag EndEvents
#tag Events sldScalAPSwitch
	#tag Event
		Sub ValueChanged()
		  ScalCalc
		  
		End Sub
	#tag EndEvent
#tag EndEvents
#tag ViewBehavior
	#tag ViewProperty
		Name="Name"
		Visible=true
		Group="ID"
		Type="String"
		EditorType="String"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Interfaces"
		Visible=true
		Group="ID"
		Type="String"
		EditorType="String"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Super"
		Visible=true
		Group="ID"
		Type="String"
		EditorType="String"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Placement"
		Visible=true
		Group="Position"
		InitialValue="0"
		Type="Integer"
		EditorType="Enum"
		#tag EnumValues
			"0 - Default"
			"1 - Parent Window"
			"2 - Main Screen"
			"3 - Parent Window Screen"
			"4 - Stagger"
		#tag EndEnumValues
	#tag EndViewProperty
	#tag ViewProperty
		Name="Width"
		Visible=true
		Group="Position"
		InitialValue="600"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Height"
		Visible=true
		Group="Position"
		InitialValue="400"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MinWidth"
		Visible=true
		Group="Position"
		InitialValue="64"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MinHeight"
		Visible=true
		Group="Position"
		InitialValue="64"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MaxWidth"
		Visible=true
		Group="Position"
		InitialValue="32000"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MaxHeight"
		Visible=true
		Group="Position"
		InitialValue="32000"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Frame"
		Visible=true
		Group="Appearance"
		InitialValue="0"
		Type="Integer"
		EditorType="Enum"
		#tag EnumValues
			"0 - Document"
			"1 - Movable Modal"
			"2 - Modal Dialog"
			"3 - Floating Window"
			"4 - Plain Box"
			"5 - Shadowed Box"
			"6 - Rounded Window"
			"7 - Global Floating Window"
			"8 - Sheet Window"
			"9 - Metal Window"
			"11 - Modeless Dialog"
		#tag EndEnumValues
	#tag EndViewProperty
	#tag ViewProperty
		Name="Composite"
		Visible=true
		Group="Appearance"
		InitialValue="False"
		Type="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="HasBackColor"
		Visible=true
		Group="Appearance"
		InitialValue="False"
		Type="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="BackColor"
		Visible=true
		Group="Appearance"
		InitialValue="&hFFFFFF"
		Type="Color"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Backdrop"
		Visible=true
		Group="Appearance"
		Type="Picture"
		EditorType="Picture"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Title"
		Visible=true
		Group="Appearance"
		InitialValue="Untitled"
		Type="String"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Visible"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="FullScreen"
		Group="Appearance"
		InitialValue="False"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MenuBarVisible"
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="CloseButton"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="Resizeable"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="LiveResize"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MaximizeButton"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MinimizeButton"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="FullScreenButton"
		Visible=true
		Group="Appearance"
		InitialValue="False"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MacProcID"
		Visible=true
		Group="Appearance"
		InitialValue="0"
		Type="Integer"
	#tag EndViewProperty
	#tag ViewProperty
		Name="MenuBar"
		Visible=true
		Group="Appearance"
		Type="MenuBar"
		EditorType="MenuBar"
	#tag EndViewProperty
	#tag ViewProperty
		Name="ImplicitInstance"
		Visible=true
		Group="Appearance"
		InitialValue="True"
		Type="Boolean"
		EditorType="Boolean"
	#tag EndViewProperty
#tag EndViewBehavior
